package spg

import (
	"math"
	"math/big"

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

/**
 ** Copyright 2018 AgileBits, Inc.
 ** Licensed under the Apache License, Version 2.0 (the "License").
 **/
