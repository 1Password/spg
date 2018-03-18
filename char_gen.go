package spg

import (
	"fmt"
	"sort"
	"strings"
)

// Character types for Character and Separator generation
const ( // character types
	ctUpper      = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	ctLower      = "abcdefghijklmnopqrstuvwxyz"
	ctDigits     = "0123456789"
	ctAmbiguous  = "0O1Il5S"
	ctSymbols    = "!#%)*+,-.:=>?@]^_}~"
	ctWhiteSpace = " \t"
)

// CTFlag is the type for the be
type CTFlag uint32

// Character type flags
const (
	Uppers CTFlag = 1 << iota
	Lowers
	Digits
	Symbols
	Ambiguous
	WhiteSpace

	None    CTFlag = 0
	Letters        = Uppers | Lowers
	All            = Letters | Digits | Symbols // This is not really all, but it is all the sane ones
)

// charTypesByFlag
var charTypeByFlag = map[CTFlag]string{
	Uppers:     ctUpper,
	Lowers:     ctLower,
	Digits:     ctDigits,
	Symbols:    ctSymbols,
	Ambiguous:  ctAmbiguous,
	WhiteSpace: ctWhiteSpace,
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
		c := chars[Int31n(uint32(len(chars)))]
		tokens[i] = Token{c, AtomTokenType}
	}
	p.Tokens = tokens
	p.Entropy = r.Entropy()
	return p, nil
}

// buildCharacterList constructs the "alphabet" that is all and only those
// characters (actually strings of length 1) that are all and only those
// characters from which the password will be build. It also ensures that
// there are no duplicates
func (r CharRecipe) buildCharacterList() []string {

	ab := r.IncludeExtra
	exclude := r.ExcludeExtra
	for f, ct := range charTypeByFlag {
		if r.Allow&f != 0 {
			ab += ct
		}
		// Treat Require as Allow for now
		if r.Require&f != 0 {
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
type CharRecipe struct {
	Length       int    // Length of generated password in characters
	Allow        CTFlag // Flags for which character types to allow
	Require      CTFlag // Flags for which character types to require
	Exclude      CTFlag // Flags for which character types to exclude
	ExcludeExtra string // Specific characters caller may want excluded
	IncludeExtra string // Specific characters caller may want excluded (this is where to put emojis. Please don't)
}

// NewCharRecipe creates CharRecipe with reasonable defaults and Length length
// Defaults are
//    r.Allow = Letters | Digits | Symbols
//    r.Exclude = Ambiguous
// And these may need to be cleared if you want to tinker with them
// This function exists only as a parallel to NewWLRecipe. It probably makes sense
// for users to forego this function and just use r := &CharRecipe{...} instead.
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
