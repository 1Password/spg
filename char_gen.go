package spg

import (
	"fmt"
	"sort"
	"strings"
)

// Character types for Character and Separator generation
const ( // character types
	ctUpper     = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	ctLower     = "abcdefghijklmnopqrstuvwxyz"
	ctDigits    = "0123456789"
	ctAmbiguous = "0O1Il5S"
	ctSymbols   = "!#%)*+,-.:=>?@]^_}~"
)

// CTFlag is the type for the be
type CTFlag uint32

// Character type flags
const (
	// Character types useful for Allow and Include
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

/*** Character type passwords ***/

// Generate a password using the character generator. The attributes contain
// all of the details needed for generating the password
func (r CharRecipe) Generate() (*Password, error) {

	if r.Length < 1 {
		return nil, fmt.Errorf("don't ask for passwords of length %d", r.Length)
	}

	p := &Password{}
	chars := r.buildCharacterList()

	tokens := make([]Token, r.Length)
	for i := 0; i < r.Length; i++ {
		c := chars[int31n(uint32(len(chars)))]
		tokens[i] = Token{c, AtomType}
	}
	p.tokens = tokens
	p.Entropy = r.Entropy()
	return p, nil
}

// buildCharacterList constructs the "alphabet" that is all and only those
// characters (actually strings of length 1) that are all and only those
// characters from which the password will be build. It also ensures that
// there are no duplicates
func (r CharRecipe) buildCharacterList() []string {

	ab := r.AllowChars
	exclude := r.ExcludeChars
	for f, ct := range charTypeByFlag {
		if r.Allow&f != 0 {
			ab += ct
		}
		// Treat Include as Allow for now
		if r.Include&f != 0 {
			ab += ct
		}
		if r.Exclude&f != 0 {
			exclude += ct
		}
	}

	alphabet := subtractString(ab, exclude)
	return strings.Split(alphabet, "")
}

// Entropy returns the entropy of a character password given the generator attributes
func (r CharRecipe) Entropy() float32 {
	size := len(r.buildCharacterList())
	return float32(entropySimple(r.Length, size))
}

// CharRecipe are generator attributes relevent for character list generation
//
// Allow - Any character from any of these sets may be present in generated password.
//
// Exclude - No characters from any of these sets may be present in the generated password.
// Exclusion overrides Include and Allow.
//
// Include - At least one character from each of these sets must be present in the generated password.
type CharRecipe struct {
	Length int // Length of generated password in characters

	// Character types to Allow, Include (require), or Exclude in generated password
	Allow   CTFlag // Types which may appear
	Include CTFlag // Types which must appear (at least one from each type)
	Exclude CTFlag // Types must not appear

	// User provided character sets for Allow, Include, and Exclude
	AllowChars   string   // Specific characters that may appear
	IncludeSets  []string // Not yet implemented
	ExcludeChars string   // Specific characters that must not appear
}

// NewCharRecipe creates CharRecipe with reasonable defaults and Length length
// Defaults are
//    r.Allow = Letters | Digits | Symbols
//    r.Exclude = Ambiguous
// And these may need to be cleared if you want to tinker with them
func NewCharRecipe(length int) *CharRecipe {

	r := new(CharRecipe)
	r.Length = length

	r.Allow = Letters | Digits | Symbols
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
