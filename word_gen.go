package spg

import (
	"fmt"
	"math"
	"strings"
)

// WLRecipe (Word List password Attributes) are the generator settings for wordlist (syllable list) passwords
type WLRecipe struct {
	list          *WordList  // Set of words for generating passwords
	Length        int        // Length of generated password in words
	SeparatorChar string     // What character(s) should separate words
	SeparatorFunc SFFunction // function to generate separators, If nil just use SeperatorChar
	Capitalize    CapScheme  // Which words in generated password should be capitalized
}

// CapScheme is for an enumeration of capitalization schemes
type CapScheme string

// Defined capitalization schemes. (Using strings instead of int enum
// to make life easier in a debugger and calling from JavaScript)
const (
	CSNone   CapScheme = "none"   // No words will be capitalized
	CSFirst  CapScheme = "first"  // First word will be capitalized
	CSAll    CapScheme = "all"    // All words will be capitalized
	CSRandom CapScheme = "random" // Some words (roughly half) will be capitalized
	CSOne    CapScheme = "one"    // One randomly selected word will be capitalized
)

// NewWLRecipe sets up word list password attributes with defaults and Length length
func NewWLRecipe(length int, wl *WordList) *WLRecipe {
	attrs := &WLRecipe{
		Length:     length,
		Capitalize: CSNone,
		list:       wl,
	}
	return attrs
}

// WordList contains the list of words WLGenerator()
type WordList struct {
	words                []string
	unCapitalizableCount int
}

// Size of the wordlist in the recipe
func (r WLRecipe) Size() uint32 {
	return r.list.Size()
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

// NewWordList does what it says on the tin. Pass it a slice of strings
// It will remove duplicates from the slice provided, and it
// will count up how many words on the list can be changed through capitalization
// This isn't cheap, so it is best to create each word list once and keep it around
// as long as you need it.
func NewWordList(list []string) (*WordList, error) {
	if len(list) == 0 {
		return nil, fmt.Errorf("cannot set up word list generator without words")
	}

	// Our RNG for picking from a list returns a uint32, so that places an upper limit on size of list
	if uint64(len(list)) > uint64(math.MaxUint32) {
		return nil, fmt.Errorf("we can't handle more than %d words", uint32(0xFFFFFFFF))
	}

	// We want to ensure that no item appears more than once
	unique := make(map[string]bool)

	for _, word := range list {
		if !unique[word] {
			unique[word] = true
		}
	}

	// A second pass to find out how many words have distinct capitalizations
	// This also treats "Polish" and "polish" as duplicates, and will
	// remove the Capitalized one from the list
	//
	// This pass also assumes that everything in unique is "true"
	unCapable := 0
	for w := range unique {
		if unique[w] { // it may have been deleted since range was computed
			cap := strings.Title(w)
			if unique[cap] {
				if cap != w { // w is "polish"
					delete(unique, cap) // delete won't change what is in range
				} else {
					unCapable++
				}
			}
		}
	}

	// third pass, because life sucks
	var ourWords []string
	for w := range unique {
		ourWords = append(ourWords, w)

	}

	if len(list) > len(ourWords) {
		// We just need to log a warning here. Not sure how we are handling that.
		// I could create a brain with standard logger and use that, but that seems
		// wrong. So let's just do this
		fmt.Printf("%d duplicate words found when setting up word list generator\n", len(list)-len(ourWords))
	}
	result := &WordList{
		words:                ourWords,
		unCapitalizableCount: unCapable,
	}
	return result, nil
}

// Generate a password using the wordlist recipe.
func (r WLRecipe) Generate() (*Password, error) {
	p := &Password{}

	if r.Size() == 0 {
		return nil, fmt.Errorf("wordlist generator must be set up before being used")
	}
	if r.Length < 1 {
		return nil, fmt.Errorf("don't ask for passwords of length %d", r.Length)
	}

	sf := r.SF()
	sepP, err := sf(r.Length - 1)
	if err != nil {
		return nil, fmt.Errorf("couldn't generate separators: %v", err)
	}

	// Construct a map of which words to capitalize
	capWords := make(map[int]bool, r.Length)
	switch r.Capitalize {
	case CSFirst:
		capWords[0] = true
	case CSOne:
		w := int(randomUint32n(uint32(r.Length)))
		capWords[w] = true
	case CSRandom:
		for i := 0; i < r.Length; i++ {
			if randomUint32n(2) == 1 {
				capWords[i] = true
			}
		}
	case CSAll:
		for i := 0; i < r.Length; i++ {
			capWords[i] = true
		}
	}

	ts := []Token{}
	for i := 0; i < r.Length; i++ {
		w := r.list.words[randomUint32n(uint32(r.Size()))]

		if capWords[i] {
			w = strings.Title(w)
		}
		if len(w) > 0 {
			ts = append(ts, Token{w, AtomType})
		}
		if i < r.Length-1 {
			sep := sepP.tokens[i].value
			if len(sep) > 0 {
				ts = append(ts, Token{sep, SeparatorType})
			}
		}
	}
	p.tokens = ts
	p.Entropy = r.Entropy()
	return p, nil
}

// Entropy returns the min-entropy from the recipe. It needs to know things
// about the wordlist used as well as other details of the recipe.
//
// When the generator produces uniform distirbution (the typical case) min-entropy
// and Shannon entropy are the same. If capitalization is used and the word list
// contains members whose capitalization does not yield a distinct element,
// the distribution becomes non-uniform.
func (r WLRecipe) Entropy() float32 {
	size := int(r.Size())
	ent := entropySimple(r.Length, size)

	// Contribution of Capitalization scheme
	if r.list.isAllCapitalizable() {
		switch r.Capitalize {
		case CSRandom:
			ent += FloatE(float64(r.Length))
		case CSOne:
			ent += FloatE(math.Log2(float64(r.Length)))
		default: // No change in entropy
		}
	}
	// else there is no additional entropy contribution from capitalization

	// Entropy contribution of separators
	sf := r.SF()
	sp, _ := sf(r.Length - 1)
	ent += FloatE(sp.Entropy)

	return float32(ent)
}

func (wl *WordList) isAllCapitalizable() bool {
	if wl.unCapitalizableCount > 0 {
		return false
	}
	return true
}

func (wl *WordList) capitalizeRatio() float64 {
	s := float64(len(wl.words))
	return (s - float64(wl.unCapitalizableCount)) / s
}

// SF will construct a separator function from the recipe
func (r WLRecipe) SF() SFFunction {
	if r.SeparatorFunc != nil {
		return r.SeparatorFunc
	}
	if len(r.SeparatorChar) > 0 {
		return SFConstant(r.SeparatorChar)
	}
	return SFConstant("")
}

/**
 ** Copyright 2018, 2019 AgileBits, Inc.
 ** Licensed under the Apache License, Version 2.0 (the "License").
 **/
