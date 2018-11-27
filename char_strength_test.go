package spg

import (
	"math"
	"testing"
)

// expectation is data for a single entropy test
type expectation struct {
	Length int

	Allow   CTFlag
	Require CTFlag
	Exclude CTFlag

	AllowChars   string
	RequireSets  []string
	ExcludeChars string

	N int64
}

const upper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const lower = "abcdefghijklmnopqrstuvwxyz"
const digits = "0123456789"

var expectations = []expectation{
	{RequireSets: []string{""}, Length: 1, N: 0},
	{RequireSets: []string{"a"}, Length: 0, N: 0},
	{RequireSets: []string{"a"}, Length: 1, N: 1},
	{RequireSets: []string{"a"}, Length: 5, N: 1},
	{RequireSets: []string{"abcde"}, Length: 1, N: 5},
	{RequireSets: []string{"abcde"}, Length: 2, N: 25},
	{RequireSets: []string{"a", "1"}, Length: 2, N: 2},
	{RequireSets: []string{"ab", "123"}, Length: 2, N: 12},

	{Require: Letters | Digits, Length: 0, N: 0},
	{Require: Letters | Digits, Length: 1, N: 0},
	{Require: Letters | Digits, Length: 2, N: 0},
	{Require: Letters | Digits, Length: 3, N: 40560},
	{RequireSets: []string{upper + lower, digits}, Length: 3, N: 96720},

	{Require: Uppers, Length: 3, N: 17576},
	{Allow: Uppers, Length: 3, N: 17576},
	{RequireSets: []string{upper + lower + digits}, Length: 3, N: 238328},
	{Allow: Letters | Digits, Length: 3, N: 238328},

	{RequireSets: []string{"a", "1"}, Length: 1, N: 0},
	{RequireSets: []string{"a"}, AllowChars: "1", Length: 1, N: 1},
	{RequireSets: []string{"a"}, AllowChars: "A1", Length: 2, N: 5},
	{RequireSets: []string{"a", "A"}, AllowChars: "1", Length: 2, N: 2},
	{RequireSets: []string{"a", "A"}, AllowChars: "1", Length: 3, N: 12},

	// Test excluded characters
	{Allow: Uppers, ExcludeChars: "ABC", Length: 1, N: 23},
	{Require: Uppers, ExcludeChars: "ABC", Length: 1, N: 23},
	{RequireSets: []string{"a", "AB"}, AllowChars: "12", ExcludeChars: "B2", Length: 3, N: 12},
}

func TestN(t *testing.T) {
	for i, exp := range expectations {
		recipe := &CharRecipe{
			Length:       exp.Length,
			Allow:        exp.Allow,
			Require:      exp.Require,
			Exclude:      exp.Exclude,
			AllowChars:   exp.AllowChars,
			RequireSets:  exp.RequireSets,
			ExcludeChars: exp.ExcludeChars,
		}
		recipe.buildCharacterList()
		intResult := recipe.n().Int64()

		if intResult != exp.N {
			t.Errorf("%d: result should be %d, was %d", i, exp.N, intResult)
		}
	}
}

func TestEntropy(t *testing.T) {
	recip := &CharRecipe{
		Length:      2,
		AllowChars:  "",
		RequireSets: []string{upper, lower, digits},
	}
	e := recip.Entropy()
	if !math.IsInf(float64(e), -1) {
		t.Errorf("Length is less than number of required sets: entropy should be -Inf, was %f", e)
	}
}

// Some tests for probability of Required success

func TestSuccessProbability(t *testing.T) {
	type tvec struct {
		Length int

		Allow   CTFlag
		Require CTFlag
		Exclude CTFlag

		AllowChars   string
		RequireSets  []string
		ExcludeChars string

		P float32
	}

	tvecs := []tvec{
		{Length: 5, Allow: Letters | Digits, P: 1},
		{Length: 3, RequireSets: []string{"123", "XYZ", "abc", "+*!"}, P: 0},

		{Require: Letters | Digits, Length: 3, P: 40560.0 / 238328.0},
		{RequireSets: []string{"ab", "123"}, Length: 2, P: 12.0 / 25.0},

		// Following tests were not calculated independently.
		{RequireSets: []string{lower, upper, digits, "+-()*&.;$#@"}, Length: 8, P: 0.444025},
		{RequireSets: []string{lower, upper, digits, "+-()*&.;$#@"}, Length: 9, P: 0.521395},
		{RequireSets: []string{lower, upper, digits, "+-()*&.;$#@"}, Length: 10, P: 0.588561},
		{RequireSets: []string{lower, upper, digits, "+-()*&.;$#@"}, Length: 11, P: 0.646499},
		{RequireSets: []string{lower, upper, digits, "+-()*&.;$#@"}, Length: 12, P: 0.696338},
		{RequireSets: []string{lower, upper, digits, "+-()*&.;$#@"}, Length: 13, P: 0.739163},
		{RequireSets: []string{lower, upper, digits, "+-()*&.;$#@"}, Length: 14, P: 0.775953},
		{RequireSets: []string{lower, upper, digits, "+-()*&.;$#@"}, Length: 15, P: 0.807565},
		{RequireSets: []string{lower, upper, digits, "+-()*&.;$#@"}, Length: 16, P: 0.834729},
		{RequireSets: []string{lower, upper, digits, "+-()*&.;$#@"}, Length: 17, P: 0.858070},
		{RequireSets: []string{lower, upper, digits, "+-()*&.;$#@"}, Length: 18, P: 0.878131},
		{RequireSets: []string{lower, upper, digits, "+-()*&.;$#@"}, Length: 19, P: 0.895364},
		{RequireSets: []string{lower, upper, digits, "+-()*&.;$#@"}, Length: 20, P: 0.910173},
	}

	for i, exp := range tvecs {
		recipe := &CharRecipe{
			Length:       exp.Length,
			Allow:        exp.Allow,
			Require:      exp.Require,
			Exclude:      exp.Exclude,
			AllowChars:   exp.AllowChars,
			RequireSets:  exp.RequireSets,
			ExcludeChars: exp.ExcludeChars,
		}
		recipe.buildCharacterList()
		p := recipe.successProbability()

		if cmpFloat32(p, exp.P, 10000) != 0 {
			t.Errorf("%d: result should be %f, was %f", i, exp.P, p)
		}
	}
}

/**
 ** Copyright 2018 AgileBits, Inc.
 ** Licensed under the Apache License, Version 2.0 (the "License").
 **/
