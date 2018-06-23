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

// required is a type for the set of strings that characters are required from
type required []string

// charList is a slice of individual characters (each as a string type)
type charList []string

/*** Character type passwords ***/

// Generate a password using the character generator. The attributes contain
// all of the details needed for generating the password
func (r CharRecipe) Generate() (*Password, error) {

	if r.Length < 1 {
		return nil, fmt.Errorf("don't ask for passwords of length %d", r.Length)
	}

	p := &Password{}
	p.Entropy = r.Entropy() // does not yet deal with inclusion requirements

	chars, include := r.buildCharacterList()

	// The difficulty of meeting requirements can be partially determined from
	// the Entropy calculation, once we calculate that properly
	if r.Length < len(include) {
		return nil, fmt.Errorf("password too short to meet all inclusion requirements")
	}

	trials := 25 // We will set this more intellegently once we have math implemented
	for i := 0; i < trials; i++ {
		tokens := make([]Token, r.Length)
		for i := 0; i < r.Length; i++ {
			c := chars[int31n(uint32(len(chars)))]
			tokens[i] = Token{c, AtomType}
		}
		p.tokens = tokens

		if includeFilter(p.String(), include) {
			return p, nil
		}
	}
	return nil, fmt.Errorf("couldn't generate password complying with requirements after %v attempts", trials)
}

// incudeFilter checks whether a candidate password has a character
// from each required/include character set
func includeFilter(pwd string, include required) bool {
	for _, rset := range include {
		if !strings.ContainsAny(pwd, rset) {
			return false
		}
	}
	return true
}

// buildCharacterList constructs the "alphabet" that is all and only those
// characters (actually strings of length 1) that are all and only those
// characters from which the password will be build. It also ensures that
// there are no duplicates
func (r CharRecipe) buildCharacterList() (charList, required) {

	ab := r.AllowChars
	exclude := r.ExcludeChars
	include := r.IncludeSets
	for f, ct := range charTypeByFlag {
		if r.Allow&f != 0 {
			ab += ct
		}
		// Include automatically gets added to alphabet, and include sets
		if r.Include&f != 0 {
			include = append(include, ct)
			ab += ct
		}
		if r.Exclude&f != 0 {
			exclude += ct
		}
	}

	alphabet := subtractString(ab, exclude)
	return strings.Split(alphabet, ""), include
}

// Entropy returns the entropy of a character password given the generator attributes
func (r CharRecipe) Entropy() float32 {
	cl, _ := r.buildCharacterList()
	size := len(cl)
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
	IncludeSets  required // Partially implemented
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
	s, _ := r.buildCharacterList()
	sort.Strings(s)
	return strings.Join(s, "")
}
