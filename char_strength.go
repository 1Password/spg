package spg

import (
	"log"
	"math"
	"math/big"
	"strings"

	set "github.com/deckarep/golang-set"
)

// This is where we do the math for the entropy calculation for character
// passwords. The trick is for when a character is _required_ from a particular set

func (r CharRecipe) entropyWithRequired() float32 {
	intValue := r.n()
	floatValue := big.NewFloat(0).SetInt(intValue)

	// big.Float doesn't have a Log function, so we need to stuff the
	// result into a float64.
	// See https://github.com/golang/go/issues/14102
	float64Value, _ := floatValue.Float64()
	if math.IsInf(float64Value, 1) {
		float64Value = math.MaxFloat64
	}

	return float32(math.Log2(float64Value))
}

func (r CharRecipe) n() *big.Int {
	allowed := set.NewSet()
	allowed.Add(r.allowedSet)
	required := set.NewSet()
	for _, req := range r.requiredSets {
		required.Add(req.s)
	}

	return n(allowed, required, r.Length)
}

// n is the number of possible passwords that can be generated.
// Unfortunately, we can't take the log until the very end, so we will
// be dealing with some very large numbers.
func n(allowed set.Set, required set.Set, length int) *big.Int {
	// totalCount is the total number of permutations possible when a
	// password of length n is generated from the set R, which is the
	// union of all sets in the password recipe.
	R := unionAll(allowed.Union(required))
	totalCount := &big.Int{}
	totalCount.Exp(toBigInt(R.Cardinality()), toBigInt(length), nil)

	// Each of these sets of sets represents a password recipe that we
	// will reject and thus must subtract from our total count.
	// We want to reject all subsets of the set of required sets except
	// the set of required sets itself.
	// For example, if L and D are required, rejectedSubsets
	// will contain {L} and {D} and will not contain {L, D}.
	// Optional sets are not part of this at all because they will
	// simply be tacked on at the end.
	powerSet := required.PowerSet()
	rejectedSubsets := set.NewSet()
	for el := range powerSet.Iter() {
		elSet, ok := el.(set.Set)
		if ok && !required.Equal(elSet) {
			rejectedSubsets.Add(elSet)
		}
	}

	// When requiredSets is {{}} (it is a set containing only the empty set),
	// powerSet(requiredSets) will also be {{}};
	// thus, rejectedSubsets will be empty, the reducing
	// function below will not run, and rejectedCount will be 0,
	// terminating the recursion.

	rejectedCount := sumAll(
		rejectedSubsets,
		func(subset set.Set) *big.Int {
			return n(allowed, subset, length)
		},
	)

	return totalCount.Sub(totalCount, rejectedCount)
}

func toBigInt(i int) *big.Int {
	return big.NewInt(int64(i))
}

// Mimic the mathematical sum operator
func sumAll(s set.Set, transform func(s set.Set) *big.Int) *big.Int {
	sum := &big.Int{}
	for el := range s.Iter() {
		elSet, ok := el.(set.Set)
		if ok {
			sum.Add(sum, transform(elSet))
		}
	}
	return sum
}

// Mimic the mathematical big union (bigcup) operator
func unionAll(elements set.Set) set.Set {
	combined := set.NewSet()
	for el := range elements.Iter() {
		elSet, ok := el.(set.Set)
		if ok {
			combined = combined.Union(elSet)
		}
	}
	return combined
}

// successProbability returns the chances of meeting all of the Require-ments
// on a single trial.
func (r CharRecipe) successProbability() float32 {

	/* The probability of generating a password that meets the Requirements
	   on a single trial is the ratio of r.n()/rWithAllRequiredChangedToAllowed.n()

	   But to avoid having to read the Go docs about big Quotients, replace that
	   division with a substraction of their logarithms. Conveniently, we have
	   those as the Entropy. Then we just raise 2 to that difference.
	*/

	newAllowChars := r.AllowChars + strings.Join(r.RequireSets, "")
	newRequireSets := []string{}
	rAllow := r
	rAllow.AllowChars = newAllowChars
	rAllow.RequireSets = newRequireSets
	rAllow.Allow = r.Allow | r.Require
	rAllow.Require = None

	eDiff := rAllow.Entropy() - r.Entropy()
	if eDiff > 0.0 {
		// This should never happen, but I don't want to
		// log.Fatal in a library
		log.Println("successProbability: eDiff is positive. Setting to 0")
		eDiff = 0.0
	}

	p := float32(math.Exp2(float64(eDiff)))
	if p > 1.0 {
		// Can't happen, but still
		log.Println("successProbability: p greater than 1. Setting to 1")
		p = 1.0
	}
	return p
}

/**
 ** Copyright 2018 AgileBits, Inc.
 ** Licensed under the Apache License, Version 2.0 (the "License").
 **/
