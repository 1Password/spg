package spg

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
	"testing"
)

const doFallibleTests = false
const entCompTolerance = 10000 // entropy must be correct to 1 part in ten thousand

var abWords []string                            // this is where we will put the wordlist for testing
const wordsFilePath = "testdata/AgileWords.txt" // relative to where test in invoked

var abSyllables []string                               // this is where we will put the wordlist for testing
const syllableFilePath = "testdata/AgileSyllables.txt" // relative to where test in invoked

func init() {
	wlFile, err := os.Open(wordsFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer wlFile.Close()

	scanner := bufio.NewScanner(wlFile)
	for scanner.Scan() {
		abWords = append(abWords, string(scanner.Text()))
	}

	slFile, err := os.Open(syllableFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer slFile.Close()

	scanner = bufio.NewScanner(slFile)
	for scanner.Scan() {
		abSyllables = append(abSyllables, string(scanner.Text()))
	}

}

func TestNewWordListPasswordGenerator(t *testing.T) {
	// First check that an empty lists returns an error
	BadG, err := NewWordListPasswordGenerator([]string{})
	if err == nil {
		t.Error("Empty wordlist should produce an error")
	}
	if BadG != nil {
		t.Error("Empty wordlist should produce a nil generator")
	}

	threeG, err := NewWordListPasswordGenerator([]string{"one", "two", "three"})
	if err != nil {
		t.Errorf("Error when creating simple wl generator: %s", err)
	}
	if threeG == nil {
		t.Error("Our three word generator should be valid, not nil")
	}

	if threeG.Size() != 3 {
		t.Errorf("there are only three words in this list, not %d", threeG.Size())
	}
}

func TestWLGenerator(t *testing.T) {

	// OK. Now for a simple wordlist test
	wordG, err := NewWordListPasswordGenerator(abWords)
	if err != nil {
		t.Errorf("Failed to create wordlist generator: %s", err)
	}
	wordAttr := NewGenAttrs(3)
	wordAttr.SeparatorChar = " "
	p, err := wordG.Generate(*wordAttr)
	pwd, ent := p.String(), p.Entropy()
	if err != nil {
		t.Errorf("failed to generate password: %s", err)
	}

	// Set up regexp under assumption that words on list are [[:alpha:]] only
	// (Sorry for all of the little pieces. I had a small error when I
	// did this all in one step)
	wRE := "\\p{L}+" // unicode letter
	sepRE := "\\Q" + wordAttr.SeparatorChar + "\\E"
	preCount := "{" + strconv.Itoa(wordAttr.Length-1) + "}"
	leadRE := "(?:" + wRE + sepRE + ")" + preCount
	res := "^" + leadRE + wRE + "$"
	re, err := regexp.Compile(res)
	if err != nil {
		t.Errorf("regexp %q did not compile: %s", re, err)
	}

	if !re.MatchString(pwd) {
		t.Errorf("pwd %q didn't match regexp %q", pwd, re)
	}
	_ = ent // keep compiler happy until I write those tests

	// As long as the test wordlist has at least 16384 the entropy for
	// for a three word password should be at least 42
	if wordG.Size() < 16384 {
		t.Errorf("this test expects a word list of at least 2^14 items. Not %d", wordG.Size())
	}
	if ent < 42.0 {
		t.Errorf("entropy (%.4f) of generated password is smaller than expected", ent)
	}

	// Let's do some math on a fixed Generator
	threeG, err := NewWordListPasswordGenerator([]string{"one", "two", "three"})
	if err != nil {
		t.Errorf("failed to create WL generator: %v", err)
	}

	p, err = threeG.Generate(GenAttrs{Length: 100})
	ent = p.Entropy()
	const expectedEnt = float32(158.496250) // 100 * log2(3). Calculated with something other than go
	if err != nil {
		t.Errorf("failed to generate long password: %v", err)
	}
	if cmpFloat32(ent, expectedEnt, entCompTolerance) != 0 {
		t.Errorf("expected entropy (%.6f) != returned entropy (%.6f)", expectedEnt, ent)
	}
}
func TestWLCapitalization(t *testing.T) {

	threeG, err := NewWordListPasswordGenerator([]string{"one", "two", "three"})
	if err != nil {
		t.Errorf("failed to create WL generator: %v", err)
	}
	// Test with random capitalization
	length := 20
	attrs := NewGenAttrs(length)
	attrs.SeparatorChar = " "
	attrs.Capitalize = CSRandom
	p, err := threeG.Generate(*attrs)
	ent := p.Entropy()
	expectedEnt := float32(51.69925) // 20 * (log2(3) + 1)
	if err != nil {
		t.Errorf("failed to generate %d word password: %v", length, err)
	}
	if cmpFloat32(ent, expectedEnt, entCompTolerance) != 0 {
		t.Errorf("expected entropy (%.6f) != returned entropy (%.6f)", expectedEnt, ent)
	}

}

func TestWLFirstCap(t *testing.T) {
	threeG, err := NewWordListPasswordGenerator([]string{"one", "two", "three"})
	if err != nil {
		t.Errorf("failed to create WL generator: %v", err)
	}
	// Test with random capitalization
	length := 5
	attrs := NewGenAttrs(length)
	attrs.SeparatorChar = " "
	attrs.Capitalize = CSFirst

	for i := 0; i < 20; i++ {
		p, err := threeG.Generate(*attrs)
		ent := p.Entropy()
		expectedEnt := float32(7.92481) // 5 * (log2(3))
		if err != nil {
			t.Errorf("failed to generate %d word password: %v", length, err)
		}
		if cmpFloat32(ent, expectedEnt, entCompTolerance) != 0 {
			t.Errorf("expected entropy (%.6f) != returned entropy (%.6f)", expectedEnt, ent)
		}
		firstWRE := "\\p{Lu}\\p{Ll}+"
		wRE := "\\p{Ll}+" // unicode letter
		sepRE := "\\Q" + attrs.SeparatorChar + "\\E"
		preCount := "{" + strconv.Itoa(attrs.Length-2) + "}"
		leadRE := firstWRE + sepRE + "(?:" + wRE + sepRE + ")" + preCount
		res := "^" + leadRE + wRE + "$"
		re, err := regexp.Compile(res)
		if err != nil {
			t.Errorf("regexp %q did not compile: %v", res, err)
		}
		if !re.MatchString(p.String()) {
			t.Errorf("%q doesn't match %s", p, re)
		}
	}

}

func TestWLRandCapitalDistribution(t *testing.T) {

	if !doFallibleTests {
		t.Skipf("Skipping statistically fallible test: %v", t.Name())
	}

	threeG, err := NewWordListPasswordGenerator([]string{"egy", "kettő", "három"})
	if err != nil {
		t.Errorf("failed to create WL generator: %v", err)
	}
	length := 1024 // big enough to make misses unlikely, round enough for me to do math easily
	attrs := NewGenAttrs(length)
	attrs.SeparatorChar = " "
	attrs.Capitalize = CSRandom
	p, _ := threeG.Generate(*attrs)
	pw := p.String()
	// We need to count the title case and non-title case words in the password
	tCaseRE, err := regexp.Compile("\\b\\p{Lu}")
	if err != nil {
		t.Errorf("title case regexp didn't compile: %v", err)
	}
	lCaseRE, err := regexp.Compile("\\b\\p{Ll}")
	if err != nil {
		t.Errorf("lowercase word regexp didn't compile: %v", err)
	}

	tCount := len(tCaseRE.FindAllString(pw, -1))
	lCount := len(lCaseRE.FindAllString(pw, -1))
	if tCount+lCount != length {
		t.Errorf("Count of title case (%d) and lower case (%d) don't add to %d", tCount, lCount, length)
	}
	if tCount < 256 || lCount < 256 {
		// chance of hitting this error by coincidence is less than 10^{-59}
		t.Errorf("far too few or too many lower case words (%d)", lCount)
	}

}
