package spg

/*** Separator functions

	Wordlist (syllable list) type generators need separators between the words,
	and creating and setting separator functions is useful. That is what is
	defined in this section.

	The over-all strategy is to generate a character type password that contains all
	of the separators needed instead of generating a separate separator for each separation.
	This allows us to use the full power of our character recipes so that we can do things
	like saying that should in toto have at least one digit and at least one symbol

***/

// A SeparatorRecipe doesn't necessarily have a length, but it may have
// a tokenizer instructions for when separator isn't just a single character
// between words
type SeparatorRecipe struct {
	cr CharRecipe
	t  *sfTokenizer // we aren't making any use of this now
}

// sfTokenizer will be instructions for how to tokenize the generated separator
// string so that its parts can be selected as needed.
// This may be useful to help specify cases where we want different separations in the same password
// to behave differently. In anticipation of issue https://github.com/1Password/spg/issues/18
type sfTokenizer struct{}

func (sr SeparatorRecipe) charRecipe(length int) *CharRecipe {
	cr := &sr.cr
	cr.Length = length
	return cr
}

// SFFunctionFull is a type for a function that returns a password
// which will be used to supply the parts for separating components
// (to be used within a password) and the entropy it contributes
// "Full" in this context means uncurried
type SFFunctionFull func(SeparatorRecipe, int) (Password, error)

// SFFunction is a curried SFFunctionFull that has already consumed
// the SeparatorRecipe
type SFFunction func(int) (*Password, error)

// NewSFFunction makes a Separator Function from a CharRecipe
func NewSFFunction(r SeparatorRecipe) SFFunction {
	return func(length int) (*Password, error) {
		cr := r.charRecipe(length)
		return cr.Generate()
	}
}

// sfConstant is for when the separator is constant
func sfConstantFull(length int, s string) (*Password, error) {
	ts := make(Tokens, length)
	for i := range ts {
		ts[i] = Token{value: s, tType: AtomType}
	}
	return &Password{Entropy: 0.0, tokens: ts}, nil
}

// SFConstant consumes the constant separator leaving a function that
// just requires a length
func SFConstant(s string) SFFunction {
	var sf SFFunction
	// Give me lambdas or give me death!
	sf = func(length int) (*Password, error) { return sfConstantFull(length, s) }
	return sf
}

// Pre-baked Separator functions
var (
	SFNone               = SFConstant("")
	SFDigits1            = NewSFFunction(SeparatorRecipe{cr: CharRecipe{Allow: Digits}})                     // Single digit separator
	SFDigitsNoAmbiguous1 = NewSFFunction(SeparatorRecipe{cr: CharRecipe{Allow: Digits, Exclude: Ambiguous}}) // Single digit, no ambiguous
	SFSymbols            = NewSFFunction(SeparatorRecipe{cr: CharRecipe{Allow: Symbols}})                    // Symbols
	SFDigitsSymbols      = NewSFFunction(SeparatorRecipe{cr: CharRecipe{Require: Symbols | Digits}})         // Symbols and digits
)

/**
 ** Copyright 2018, 2019 AgileBits, Inc.
 ** Licensed under the Apache License, Version 2.0 (the "License").
 **/
