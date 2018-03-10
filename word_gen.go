package spg

import (
	"fmt"
	"math"
	"strings"
)

// WLRecipe (Word List password Attributes) are the generator settings for wordlist (syllable list) passwords
type WLRecipe struct {
	Length        int        // Length of generated password in words
	SeparatorChar string     // What character(s) should separate words
	SeparatorFunc SFFunction // function to generate separators, If nil just use SeperatorChar
	Capitalize    CapScheme  // Which words in generated password should be capitalized
}

// CapScheme is for an enumeration of capitalization schemes
type CapScheme string

// Capitalization schemes for wordlist (and syllable lists)
// (Using strings instead of ints makes for more useful error messages)
const (
	CSNone   = "none"   // No words will be capitalized
	CSFirst  = "first"  // First word will be capitalized
	CSAll    = "all"    // All words will be capitalized
	CSRandom = "random" // Some words (roughly half) will be capitalized
	CSOne    = "one"    // One randomly selected word will be capitalized
)

// NewWLRecipe sets up word list password attributes with defaults and Length length
func NewWLRecipe(length int) *WLRecipe {
	attrs := &WLRecipe{
		Length:     length,
		Capitalize: CSNone,
	}
	return attrs
}

// WordList contains the list of words WLGenerator()
type WordList struct {
	words []string
}

// Size returns the number of items in the generator's wordlist or the maxiumum uint32, whichever is smaller
// (the restriction on size is because of the RNG we are using)
func (wl WordList) Size() uint32 {
	size := len(wl.words)

	// Why all this casting? (yes, functions not casts.) Because gopherjs won't assign
	// math.MaxUint32 to an int. It doesn't like untyped values an considers it overflow
	if uint64(size) > uint64(math.MaxUint32) {
		return uint32(math.MaxUint32)
	}
	return uint32(size)
}

// NewWordList does what is says on the tin. Pass it a slice of strings
func NewWordList(list []string) (*WordList, error) {
	if len(list) == 0 {
		return nil, fmt.Errorf("cannot set up word list generator without words")
	}

	// Our RNG for picking from a list returns a uint32, so that places an upper limit on size of list
	if uint64(len(list)) > uint64(math.MaxUint32) {
		return nil, fmt.Errorf("we can't handle more than %d words", uint32(0xFFFFFFFF))
	}

	// We want to ensure that no item appears more than once
	unique := make(map[string]bool, len(list))
	var ourWords []string // Don't create with make. We need this to start with zero length
	for _, word := range list {
		if !unique[word] {
			ourWords = append(ourWords, word)
			unique[word] = true
		}
	}
	if len(list) > len(ourWords) {
		// We just need to log a warning here. Not sure how we are handling that.
		// I could create a brain with standard logger and use that, but that seems
		// wrong. So let's just do this
		fmt.Printf("%d duplicate words found when setting up word list generator\n", len(list)-len(ourWords))
	}
	result := &WordList{
		words: ourWords,
	}
	return result, nil
}

// Generate a password using the wordlist generator. Requires that the generator already be set up
// Although we are passing a pointer to a generator, that is only to avoid some
// memory copying. This does not change g.
func (r WLRecipe) Generate(g *WordList) (Password, error) {
	p := Password{}
	if g.Size() == 0 {
		return p, fmt.Errorf("wordlist generator must be set up before being used")
	}
	if r.Length < 1 {
		return p, fmt.Errorf("don't ask for passwords of length %d", r.Length)
	}

	var sf SFFunction
	if r.SeparatorFunc == nil {
		sf = SFFunction(func() (string, float64) { return r.SeparatorChar, 0.0 })
	} else {
		sf = r.SeparatorFunc
	}

	// Construct a map of which words to capitalize
	capWords := make(map[int]bool, r.Length)
	switch r.Capitalize {
	case CSFirst:
		capWords[0] = true
	case CSOne:
		w := int(Int31n(uint32(r.Length)))
		capWords[w] = true
	case CSRandom:
		for i := 1; i <= r.Length; i++ {
			if Int31n(2) == 1 {
				capWords[i] = true
			}
		}
	case CSAll:
		for i := 1; i <= r.Length; i++ {
			capWords[i] = true
		}
	}

	toks := []Token{}
	for i := 0; i < r.Length; i++ {
		w := g.words[Int31n(uint32(g.Size()))]

		if capWords[i] {
			w = strings.Title(w)
		}
		if len(w) > 0 {
			toks = append(toks, Token{w, AtomTokenType})
		}
		if i < r.Length-1 {
			sep, _ := sf()
			if len(sep) > 0 {
				toks = append(toks, Token{sep, SeparatorTokenType})
			}
		}
	}
	p.Tokens = toks
	p.ent = r.Entropy(int(g.Size()))
	return p, nil
}

// Entropy needs to know the wordlist size to calculate entropy for some attributes
// BUG(jpg) Wordlist capitalization entropy calculation assumes that all words in list begin with a lowercase letter.
func (r WLRecipe) Entropy(listSize int) float32 {
	ent := entropySimple(r.Length, listSize)
	switch r.Capitalize {
	case CSRandom:
		ent += float64(r.Length)
	case CSOne:
		ent += math.Log2(float64(r.Length))
	default: // No change in entropy
	}

	// Entropy contribution of separators
	sepEnt := 0.0
	if r.SeparatorFunc != nil {
		_, sepEnt = r.SeparatorFunc()
	}
	ent += (float64(r.Length) - 1.0) * sepEnt

	return float32(ent)
}

/*** Separator functions

	Wordlist (syllable list) type generators need separators between the words,
	and creating and setting separator functions is useful. That is what is
	defined in this section.

***/

// SFFunction is a type for a function that returns a string
// (to be used within a password) and the entropy it contributes
type SFFunction func() (string, float64)

// Pre-baked Separator functions

// SFNone empty separator
func SFNone() (string, float64) { return "", 0.0 }

// SFDigits1 each separator is a randomly chosen digit
func SFDigits1() (string, float64) { return nFromString(CTDigits, 1) }

// SFDigits2 each separator is two randomly chosen digits
func SFDigits2() (string, float64) { return nFromString(CTDigits, 2) }

// SFDigitsNoAmbiguous1 each separator is a non-ambiguous digit
func SFDigitsNoAmbiguous1() (string, float64) {
	return nFromString(subtractString(CTDigits, CTAmbiguous), 1)
}

// SFDigitsNoAmbiguous2 each separator is a pair of randomly chosen non-ambiguous digits
func SFDigitsNoAmbiguous2() (string, float64) {
	return nFromString(subtractString(CTDigits, CTAmbiguous), 2)
}

// SFSymbols each separator is a randomly chosen symbol
func SFSymbols() (string, float64) { return nFromString(CTSymbols, 1) }

// SFDigitsSymbols each separator is a randomly chosen digit or symbol
func SFDigitsSymbols() (string, float64) { return nFromString(CTSymbols+CTDigits, 1) }
