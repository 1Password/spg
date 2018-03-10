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
func (r CharRecipe) Generate() (Password, error) {

	if r.Length < 1 {
		return Password{}, fmt.Errorf("don't ask for passwords of length %d", r.Length)
	}

	p := Password{}
	chars := r.buildCharacterList()

	toks := make([]Token, r.Length)
	for i := 0; i < r.Length; i++ {
		c := chars[Int31n(uint32(len(chars)))]
		toks[i] = Token{c, AtomTokenType}
	}
	p.Tokens = toks
	p.ent = r.Entropy()
	return p, nil
}

// buildCharacterList constructs the "alphabet" that is all and only those
// characters (actually strings of length 1) that are all and only those
// characters from which the password will be build. It also ensures that
// there are no duplicates
func (r CharRecipe) buildCharacterList() []string {
	// No letters overrides any Uppers or Lowers case settings
	if r.Letter != CIUnstated {
		r.Lowers = r.Letter
		r.Uppers = r.Letter
	}

	/* We have three steps in creating the set of characters to use
	   1. Build it up from what is allowed
	   2. Remove duplicate characters from the list
	   3. Remove exclusions

	   Steps 2 and 3 are accomplished by the subtractString() function
	*/

	ab := r.IncludeExtra
	if r.Digits == CIInclude {
		ab += CTDigits
	}
	if r.Lowers == CIInclude {
		ab += CTLower
	}
	if r.Uppers == CIInclude {
		ab += CTUpper
	}
	if r.Symbols == CIInclude {
		ab += CTSymbols
	}
	if r.Ambiguous == CIInclude {
		ab += CTAmbiguous
	}

	exclude := r.ExcludeExtra
	if r.Digits == CIExclude {
		exclude += CTDigits
	}
	if r.Lowers == CIExclude {
		exclude += CTLower
	}
	if r.Uppers == CIExclude {
		exclude += CTUpper
	}
	if r.Symbols == CIExclude {
		exclude += CTSymbols
	}
	if r.Ambiguous == CIExclude {
		exclude += CTAmbiguous
	}

	alphabet := subtractString(ab, exclude)
	return strings.Split(alphabet, "")
}

// Entropy returns the entropy of a character password given the generator attributes
func (r CharRecipe) Entropy() float32 {
	size := len(r.buildCharacterList())
	return float32(entropySimple(r.Length, size))
}

// CharInclusion holds the inclusion/exclusion value for some character class
type CharInclusion string

// CI{Included,Required,Excluded,Unstated} indicate how some class of characters (such as digts)
// are to be included (or not) in the generated password
const (
	CIInclude  = CharInclusion("included") // Included in the set of characters used by generator
	CIRequire  = CharInclusion("required") // At least one of these must be in each generated password
	CIExclude  = CharInclusion("excluded") // None of these may appear in a generated password
	CIUnstated = CharInclusion("")         // Not included by this statement, but not excluded either
)

// CharRecipe are generator attributes relevent for character list generation
type CharRecipe struct {
	Length       int           // Length of generated password in characters
	Uppers       CharInclusion // Uppercase letters, [A-Z] may be included in password
	Lowers       CharInclusion // Lowercase letters, [a-z] may be included in password
	Letter       CharInclusion // If false, overrides Lowers and Uppers setting, does nothing if true
	Digits       CharInclusion // Digits [0-9] may be included in password
	Symbols      CharInclusion // Symbols, punctuation characters may be included in password
	Ambiguous    CharInclusion // Ambiguous characters (such as "I" and "1") are to be excluded from password
	ExcludeExtra string        // Specific characters caller may want excluded
	IncludeExtra string        // Specific characters caller may want excluded (this is where to put emojis. Please don't)
}

// NewCharRecipe creates CharRecipe with reasonable defaults and Length length
// more structure
func NewCharRecipe(length int) *CharRecipe {

	r := new(CharRecipe)
	r.Length = length

	r.Ambiguous = CIExclude

	r.Digits = CIInclude
	r.Uppers = CIInclude
	r.Lowers = CIInclude
	r.Symbols = CIInclude

	return r
}
