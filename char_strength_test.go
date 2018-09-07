package spg

import (
	"math"
	"testing"

	set "github.com/deckarep/golang-set"
)

// expectation is data for a single entropy test
type expectation struct {
	Allowed  []string
	Required []string
	Length   int
	Result   int64
}

const upper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const lower = "abcdefghijklmnopqrstuvwxyz"
const digits = "0123456789"

var expectations = []expectation{
	{Required: []string{""}, Length: 1, Result: 0},
	{Required: []string{"a"}, Length: 0, Result: 0},
	{Required: []string{"a"}, Length: 1, Result: 1},
	{Required: []string{"a"}, Length: 5, Result: 1},
	{Required: []string{"abcde"}, Length: 1, Result: 5},
	{Required: []string{"abcde"}, Length: 2, Result: 25},
	{Required: []string{"a", "1"}, Length: 2, Result: 2},
	{Required: []string{"ab", "123"}, Length: 2, Result: 12},

	{Required: []string{upper, lower, digits}, Length: 0, Result: 0},
	{Required: []string{upper, lower, digits}, Length: 1, Result: 0},
	{Required: []string{upper, lower, digits}, Length: 2, Result: 0},
	{Required: []string{upper, lower, digits}, Length: 3, Result: 40560},
	{Required: []string{upper + lower, digits}, Length: 3, Result: 96720},

	{Required: []string{upper}, Length: 3, Result: 17576},
	{Allowed: []string{upper}, Length: 3, Result: 17576},
	{Required: []string{upper + lower + digits}, Length: 3, Result: 238328},
	{Allowed: []string{upper, lower, digits}, Length: 3, Result: 238328},

	{Required: []string{"a", "1"}, Length: 1, Result: 0},
	{Required: []string{"a"}, Allowed: []string{"1"}, Length: 1, Result: 1},
	{Required: []string{"a"}, Allowed: []string{"A", "1"}, Length: 2, Result: 5},
	{Required: []string{"a"}, Allowed: []string{"A1"}, Length: 2, Result: 5},
	{Required: []string{"a", "A"}, Allowed: []string{"1"}, Length: 2, Result: 2},
	{Required: []string{"a", "A"}, Allowed: []string{"1"}, Length: 3, Result: 12},
}

func toSetOfSets(arr []string) set.Set {
	ret := set.NewSet()
	for _, str := range arr {
		ret.Add(setFromString(str))
	}
	return ret
}

func TestN(t *testing.T) {
	for i, exp := range expectations {
		result := n(toSetOfSets(exp.Allowed), toSetOfSets(exp.Required), exp.Length)
		intResult := result.Int64()

		if intResult != exp.Result {
			t.Errorf("%d: result should be %d, was %d", i, exp.Result, intResult)
		}
	}
}

func TestEntropy(t *testing.T) {
	recip := &CharRecipe{Length: 2}
	recip.AllowChars = ""
	recip.RequireSets = []string{upper, lower, digits}
	e := recip.Entropy()
	if !math.IsInf(float64(e), -1) {
		t.Errorf("entropy should be -Inf, was %f", e)
	}
}
