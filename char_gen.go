package spg

import (
	"fmt"
	"strings"
)

// Character types for Character and Separator generation
const ( // character types
	CTUpper      = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	CTLower      = "abcdefghijklmnopqrstuvwxyz"
	CTDigits     = "0123456789"
	CTAmbiguous  = "0O1Il5S"
	CTSymbols    = "!#%)*+,-.:=>?@]^_}~"
	CTWhiteSpace = " \t"
)

/*** Character type passwords ***/

// CharacterPasswordGenerator generates a password from random characters
type CharacterPasswordGenerator struct {
}

// NewCharacterPasswordGenerator exists only as a parallel to NewWordListPasswordGenerator
// It doesn't do anying other than return new(CharacterPasswordGenerator), nil
func NewCharacterPasswordGenerator() (*CharacterPasswordGenerator, error) {
	return new(CharacterPasswordGenerator), nil
}

// Generate a password using the character generator. The attributes contain
// all of the details needed for generating the password
func (g CharacterPasswordGenerator) Generate(attrs GenAttrs) (Password, error) {

	if attrs.Length < 1 {
		return Password{}, fmt.Errorf("don't ask for passwords of length %d", attrs.Length)
	}

	p := Password{}
	chars := attrs.buildCharacterList()

	toks := make([]Token, attrs.Length)
	for i := 0; i < attrs.Length; i++ {
		c := chars[Int31n(uint32(len(chars)))]
		toks[i] = Token{c, AtomTokenType}
	}
	p.Tokens = toks
	p.ent = attrs.CEntropy()
	return p, nil
}

// buildCharacterList constructs the "alphabet" that is all and only those
// characters (actually strings of length 1) that are all and only those
// characters from which the password will be build. It also ensures that
// there are no duplicates
func (a CharGenAttrs) buildCharacterList() []string {
	// No letters overrides any Upper or Lower case settings
	if !a.AllowLetter {
		a.AllowLower = false
		a.AllowUpper = false
	}

	/* We have three steps in creating the set of characters to use
	   1. Build it up from what is allowed
	   2. Remove duplicate characters from the list
	   3. Remove exclusions

	   Steps 2 and 3 are accomplished by the subtractString() function
	*/

	ab := ""
	if a.AllowDigit {
		ab += CTDigits
	}
	if a.AllowLower {
		ab += CTLower
	}
	if a.AllowUpper {
		ab += CTUpper
	}
	if a.AllowSymbol {
		ab += CTSymbols
	}
	if a.AllowWhiteSpace {
		ab += CTWhiteSpace
	}
	ab += a.IncludeExtra

	exclude := a.ExcludeExtra
	if a.ExcludeAmbiguous {
		exclude += CTAmbiguous
	}

	alphabet := subtractString(ab, exclude)
	return strings.Split(alphabet, "")
}

// CEntropy returns the entropy of a character password given the generator attributes
func (a GenAttrs) CEntropy() float32 {
	size := len(a.buildCharacterList())
	return float32(entropySimple(a.Length, size))
}

// CharGenAttrs are generator attributes relevent for character list generation
type CharGenAttrs struct {
	AllowUpper       bool   // Uppercase letters, [A-Z] may be included in password
	AllowLower       bool   // Lowercase letters, [a-z] may be included in password
	AllowLetter      bool   // If false, overrides Lower and Upper setting, does nothing if true
	AllowDigit       bool   // Digits [0-9] may be included in password
	AllowSymbol      bool   // Symbols, punctuation characters may be included in password
	ExcludeAmbiguous bool   // Ambiguous characters (such as "I" and "1") are to be excluded from password
	AllowWhiteSpace  bool   // Allow space and tab in passwords (this is silly, don't set)
	ExcludeExtra     string // Specific characters caller may want excluded
	IncludeExtra     string // Specific characters caller may want excluded (this is where to put emojis. Please don't)
}
