package spg

import (
	"fmt"
	"math"
	"sort"
	"strings"

	set "github.com/deckarep/golang-set"
)

// Character types for Character and Separator generation
const ( // character types
	ctUpper     = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	ctLower     = "abcdefghijklmnopqrstuvwxyz"
	ctDigits    = "0123456789"
	ctAmbiguous = "0O1Il5S"
	ctSymbols   = "!#%)*+,-.:=>?@]^_}~"
)

/* We have three different internal representations of collections of characters:
 1. `string`
	These are just handy for any public API, but they can't be
	directly used for anything and they don't guarantee that elements
	aren't repeated. Strings are also useful in a strings.ContainsAny()
	construction we may use for filtering.

 2. `set.Set`
	These are useful for when we need to perform set operations such as set difference (which we do need for a number of different reasons)

3. `[]string` where each string in the slice is a single rune
	This is needed for when we need to select uniform random rune from
	the set. This representation is only needed for the total alphebet
	passwords are generated from.

	To (hopefully) avoid confusion with other arrays for strings,
	we have a type alias for this called `charList`.

In light of all of this, I'm going to go against some Go conventions
and name types in terms of their underlying types.
*/

// required is a type for the set of strings that characters are required from
type required []string

// charList is a slice of individual characters (each as a string type)
type charList []string

// CTFlag is the type for the be
type CTFlag uint32

// Character type flags
const (
	// Character types useful for Allow and Require
	Uppers CTFlag = 1 << iota
	Lowers
	Digits
	Symbols

	// Character types useful for Exclude
	Ambiguous

	// Named combinations
	None    CTFlag = 0
	Letters        = Uppers | Lowers
	All            = Letters | Digits | Symbols
)

// charTypesByFlag
var charTypeByFlag = map[CTFlag]string{
	Uppers:    ctUpper,
	Lowers:    ctLower,
	Digits:    ctDigits,
	Symbols:   ctSymbols,
	Ambiguous: ctAmbiguous,
}

// ctNamesByFlag
var charTypeNamesByFlag = map[CTFlag]string{
	Uppers:  "Uppers",
	Lowers:  "Lowers",
	Digits:  "Digits",
	Symbols: "Symbols",

	None:    "None",
	Letters: "Letters",
	All:     "All characters",
}

// Re-generation trials for meeting requirements
var (
	MaxTrials   = 200              // How many times we will try to generate before giving up
	MaxFailRate = 1.0 / 1000000000 // Maximum acceptable failure rate after MaxTrials
)

func (r CharRecipe) hasAcceptableFailRate() (bool, float32) {
	sp := r.SuccessProbability()
	if sp <= 0.0 {
		return false, 1.0
	}
	failP := math.Pow(1.0-float64(sp), float64(MaxTrials))
	return failP <= MaxFailRate, float32(failP)
}

/*** Character type passwords ***/

// Generate a password using the character generator. The attributes contain
// all of the details needed for generating the password
func (r CharRecipe) Generate() (*Password, error) {

	if r.Length < 1 {
		return nil, fmt.Errorf("don't ask for passwords of length %d", r.Length)
	}

	p := &Password{}
	p.Entropy = r.Entropy()

	chars := r.buildCharacterList()
	if len(chars) == 0 {
		return nil, fmt.Errorf("no characters to build pwd from")
	}

	if acceptable, failP := r.hasAcceptableFailRate(); !acceptable {
		return nil, fmt.Errorf("Chance of not generated a valid password (%v) is too high", failP)
	}
	for i := 0; i < MaxTrials; i++ {
		tokens := make([]Token, r.Length)
		for i := 0; i < r.Length; i++ {
			c := chars[randomUint32n(uint32(len(chars)))]
			tokens[i] = Token{c, AtomType}
		}
		p.tokens = tokens

		ps := p.String() // creating this variable for debugging
		if requireFilter(ps, r.requiredSets) {
			return p, nil
		}
	}
	return nil, fmt.Errorf("couldn't generate password complying with requirements after %v attempts", MaxTrials)
}

// buildCharacterList constructs the "alphabet" that is all and only those
// characters (actually strings of length 1) from which the password will be
// built. It also ensures that there are no duplicates.
func (r *CharRecipe) buildCharacterList() charList {
	allowedChars := r.AllowChars
	excludedChars := r.ExcludeChars
	r.requiredSets = make(reqSets, 0)
	for i, s := range r.RequireSets {
		if len(s) > 0 {
			r.requiredSets = append(r.requiredSets,
				*newReqSet(s, fmt.Sprintf("Custom %d", i+1)))
		}
	}
	for f, ct := range charTypeByFlag {
		if r.Allow&f != 0 {
			allowedChars += ct
		}
		if r.Require&f != 0 {
			ctName, ok := charTypeNamesByFlag[f]
			if !ok {
				ctName = "Dunno"
			}
			r.requiredSets = append(r.requiredSets, *newReqSet(ct, ctName))
		}
		if r.Exclude&f != 0 {
			excludedChars += ct
		}
	}

	// Now we need to clean this all up. First let's make them sets.
	excludedSet := setFromString(excludedChars)
	r.allowedSet = setFromString(allowedChars).Difference(excludedSet)

	// Now remove excluded chars from each required set
	// and remove required chars from the allowed set
	for i := range r.requiredSets {
		req := &r.requiredSets[i]
		req.s = req.s.Difference(excludedSet)
		r.allowedSet = r.allowedSet.Difference(req.s)
	}

	alphabetSet := r.allowedSet.Union(r.requiredSets.union().s)
	return strings.Split(stringFromSet(alphabetSet), "")
}

// Entropy returns the entropy of a character password given the generator attributes
func (r CharRecipe) Entropy() float32 {
	cl := r.buildCharacterList()
	if r.requiredSets.size() != 0 {
		return r.entropyWithRequired()
	}
	size := len(cl)
	return float32(entropySimple(r.Length, size))
}

// CharRecipe are generator attributes relevent for character list generation
//
// Allow - Any character from any of these sets may be present in generated password.
//
// Exclude - No characters from any of these sets may be present in the generated password.
// Exclusion overrides Require and Allow.
//
// Require - At least one character from each of these sets must be present in the generated password.
type CharRecipe struct {
	Length int // Length of generated password in characters

	// Character types to Allow, Require, or Exclude in generated password
	Allow   CTFlag // Types which may appear
	Require CTFlag // Types which must appear (at least one from each type)
	Exclude CTFlag // Types must not appear

	// User provided character sets for Allow, Require, and Exclude
	AllowChars   string   // Specific characters that may appear
	RequireSets  []string // At least one character from each string must appear
	ExcludeChars string   // Specific characters that must not appear

	// Following sets are computed
	allowedSet   set.Set // Allowed, but not required
	requiredSets reqSets // List of sets of required characters
}

// NewCharRecipe creates CharRecipe with reasonable defaults and Length length
// Defaults are
//    r.Require = Uppers | Lowers | Digits
//    r.Exclude = Ambiguous
// And these may need to be cleared if you want to tinker with them
func NewCharRecipe(length int) *CharRecipe {

	r := new(CharRecipe)
	r.Length = length

	r.Require = Uppers | Lowers | Digits
	r.Exclude = Ambiguous

	return r
}

// Alphabet returns a sorted string of the characters that are
// drawn from in a given recipe, r
func (r CharRecipe) Alphabet() string {
	s := r.buildCharacterList()
	sort.Strings(s)
	return strings.Join(s, "")
}

/**
 ** Copyright 2018 AgileBits, Inc.
 ** Licensed under the Apache License, Version 2.0 (the "License").
 **/
