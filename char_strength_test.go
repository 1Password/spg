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
		{RequireSets: []string{lower, upper, digits, ctSymbols}, Length: 1, P: 0.0},
		{RequireSets: []string{lower, upper, digits, ctSymbols}, Length: 2, P: 0.0},
		{RequireSets: []string{lower, upper, digits, ctSymbols}, Length: 3, P: 0.0},
		{RequireSets: []string{lower, upper, digits, ctSymbols}, Length: 4, P: 0.071610},
		{RequireSets: []string{lower, upper, digits, ctSymbols}, Length: 5, P: 0.179024},
		{RequireSets: []string{lower, upper, digits, ctSymbols}, Length: 6, P: 0.293271},
		{RequireSets: []string{lower, upper, digits, ctSymbols}, Length: 7, P: 0.399864},
		{RequireSets: []string{lower, upper, digits, ctSymbols}, Length: 8, P: 0.493387},
		{RequireSets: []string{lower, upper, digits, ctSymbols}, Length: 9, P: 0.572959},
		{RequireSets: []string{lower, upper, digits, ctSymbols}, Length: 10, P: 0.639650},
		{RequireSets: []string{lower, upper, digits, ctSymbols}, Length: 11, P: 0.695186},
		{RequireSets: []string{lower, upper, digits, ctSymbols}, Length: 12, P: 0.741371},
		{RequireSets: []string{lower, upper, digits, ctSymbols}, Length: 13, P: 0.779832},
		{RequireSets: []string{lower, upper, digits, ctSymbols}, Length: 14, P: 0.811963},
		{RequireSets: []string{lower, upper, digits, ctSymbols}, Length: 15, P: 0.838902},
		{RequireSets: []string{lower, upper, digits, ctSymbols}, Length: 16, P: 0.861589},
		{RequireSets: []string{lower, upper, digits, ctSymbols}, Length: 17, P: 0.880777},
		{RequireSets: []string{lower, upper, digits, ctSymbols}, Length: 18, P: 0.897070},
		{RequireSets: []string{lower, upper, digits, ctSymbols}, Length: 19, P: 0.910953},
		{RequireSets: []string{lower, upper, digits, ctSymbols}, Length: 20, P: 0.922832},
	}

	for _, exp := range tvecs {
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
		p := recipe.SuccessProbability()

		if cmpFloat32(p, exp.P, 10000) != 0 {
			t.Errorf("Result for length %d should be %f, was %f", exp.Length, exp.P, p)
		}

		// I want to call the generator on some of these in the debugger, so ...
		pwd, err := recipe.Generate()
		_ = err
		_ = pwd
	}
}

/**
 ** Copyright 2018 AgileBits, Inc.
 ** Licensed under the Apache License, Version 2.0 (the "License").
 **/
