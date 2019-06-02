package spg

/*** Separator functions

	Wordlist (syllable list) type generators need separators between the words,
	and creating and setting separator functions is useful. That is what is
	defined in this section.

***/

// A SeparatorRecipe doesn't necessarily have a length, but it may have
// a tokenizer instructions for when separator isn't just a single character
// between words
type SeparatorRecipe struct {
	cr CharRecipe
	t  *sfTokenizer
}

// sfTokenizer will be instructions for how to tokenize the generated separator
// string so that its parts can be selected as needed
type sfTokenizer struct{}

func (sr SeparatorRecipe) charRecipe(length int) *CharRecipe {
	cr := &sr.cr
	cr.Length = length
	return cr
}

// SFFunctionFull is a type for a function that returns a password
// which will be used to supply the parts for separating components
// (to be used within a password) and the entropy it contributes
type SFFunctionFull func(SeparatorRecipe, int) (Password, error)

// SFFunction is a curried SFFunctionFull, but has already consumed
// the SeparatorRecipe
type SFFunction func(int) (*Password, error)

// NewSFFunction makes a Separator Function from a CharRecipe
func NewSFFunction(r SeparatorRecipe) SFFunction {
	var sf SFFunction
	sf = func(length int) (*Password, error) { return sfWrap(r, length) }
	return sf
}

func sfWrap(sr SeparatorRecipe, length int) (*Password, error) {
	r := sr.charRecipe(length)
	return r.Generate()
}

var nullToken = Token{
	value: "",
	tType: AtomType,
}

func sfConstantFull(length int, s string) (*Password, error) {
	ts := make(Tokens, length)
	for i := range ts {
		ts[i] = Token{value: s, tType: AtomType}
	}
	return &Password{Entropy: 0.0, tokens: ts}, nil
}

func sfConstant(s string) SFFunction {
	var sf SFFunction
	sf = func(length int) (*Password, error) { return sfConstantFull(length, s) }
	return sf

}

// sfNull generates a separator password of length length with empty tokens
func sfNull(length int) (*Password, error) {
	ts := make(Tokens, length)
	for i := range ts {
		ts[i] = nullToken
	}
	return &Password{Entropy: 0.0, tokens: ts}, nil
}

// Pre-baked Separator functions
var (
	SFNone               = sfConstant("")
	SFDigits1            = NewSFFunction(SeparatorRecipe{cr: CharRecipe{Allow: Digits}})                     // Single digit separator
	SFDigitsNoAmbiguous1 = NewSFFunction(SeparatorRecipe{cr: CharRecipe{Allow: Digits, Exclude: Ambiguous}}) // Single digit, no ambiguous
	SFSymbols            = NewSFFunction(SeparatorRecipe{cr: CharRecipe{Allow: Symbols}})                    // Symbols
	SFDigitsSymbols      = NewSFFunction(SeparatorRecipe{cr: CharRecipe{Allow: Symbols | Digits}})           // Symbols and digits
)

/**
 ** Copyright 2018, 2019 AgileBits, Inc.
 ** Licensed under the Apache License, Version 2.0 (the "License").
 **/
