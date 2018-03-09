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

// Generate a password using the character generator. The attributes contain
// all of the details needed for generating the password
func (a CharGenAttrs) Generate() (Password, error) {

	if a.Length < 1 {
		return Password{}, fmt.Errorf("don't ask for passwords of length %d", a.Length)
	}

	p := Password{}
	chars := a.buildCharacterList()

	toks := make([]Token, a.Length)
	for i := 0; i < a.Length; i++ {
		c := chars[Int31n(uint32(len(chars)))]
		toks[i] = Token{c, AtomTokenType}
	}
	p.Tokens = toks
	p.ent = a.Entropy()
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

// Entropy returns the entropy of a character password given the generator attributes
func (a CharGenAttrs) Entropy() float32 {
	size := len(a.buildCharacterList())
	return float32(entropySimple(a.Length, size))
}

// CharGenAttrs are generator attributes relevent for character list generation
type CharGenAttrs struct {
	Length           int    // Length of generated password in characters
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

// NewCharAttrs creates CharGenAttrs with reasonable defaults and Length length
// more structure
func NewCharAttrs(length int) *CharGenAttrs {
	const (
		defaultSep        = ""
		defaultDigits     = true
		defaultUpper      = true
		defaultLower      = true
		defaultSymbol     = true
		defaultAmbiguous  = true // exclude ambiguous by default
		defaultWhiteSpace = false
		defaultExclude    = ""
	)
	// function literal cannot be a string

	attrs := new(CharGenAttrs)
	attrs.Length = length

	attrs.ExcludeAmbiguous = defaultAmbiguous
	attrs.ExcludeExtra = defaultExclude

	attrs.AllowDigit = defaultDigits
	attrs.AllowUpper = defaultUpper
	attrs.AllowLower = defaultLower
	attrs.AllowLetter = attrs.AllowUpper || attrs.AllowLower
	attrs.AllowSymbol = defaultSymbol
	attrs.AllowWhiteSpace = defaultWhiteSpace
	attrs.IncludeExtra = ""

	return attrs
}
