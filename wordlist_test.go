package spg

import (
	"bufio"
	"log"
	"math"
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
	BadG, err := NewWordList([]string{})
	if err == nil {
		t.Error("Empty wordlist should produce an error")
	}
	if BadG != nil {
		t.Error("Empty wordlist should produce a nil generator")
	}

	threeG, err := NewWordList([]string{"one", "two", "three"})
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
	wl, err := NewWordList(abWords)
	if err != nil {
		t.Errorf("Failed to create wordlist generator: %s", err)
	}
	a := NewWLRecipe(3, wl)
	a.SeparatorChar = " "
	p, err := a.Generate()
	pwd, ent := p.String(), p.Entropy()
	if err != nil {
		t.Errorf("failed to generate password: %s", err)
	}

	// Set up regexp under assumption that words on list are [[:alpha:]] only
	// (Sorry for all of the little pieces. I had a small error when I
	// did this all in one step)
	wRE := "\\p{L}+" // unicode letter
	sepRE := "\\Q" + a.SeparatorChar + "\\E"
	preCount := "{" + strconv.Itoa(a.Length-1) + "}"
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
	if wl.Size() < 16384 {
		t.Errorf("this test expects a word list of at least 2^14 items. Not %d", wl.Size())
	}
	if ent < 42.0 {
		t.Errorf("entropy (%.4f) of generated password is smaller than expected", ent)
	}

	// Let's do some math on a fixed Generator
	threeG, err := NewWordList([]string{"one", "two", "three"})
	if err != nil {
		t.Errorf("failed to create WL generator: %v", err)
	}

	p, err = WLRecipe{Length: 100, WordList: threeG}.Generate()
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

	threeG, err := NewWordList([]string{"one", "two", "three"})
	if err != nil {
		t.Errorf("failed to create WL generator: %v", err)
	}
	// Test with random capitalization
	length := 20
	r := NewWLRecipe(length, threeG)
	r.SeparatorChar = " "
	r.Capitalize = CSRandom
	p, err := r.Generate()
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
	threeG, err := NewWordList([]string{"one", "two", "three"})
	if err != nil {
		t.Errorf("failed to create WL generator: %v", err)
	}
	// Test with random capitalization
	length := 5
	r := NewWLRecipe(length, threeG)
	r.SeparatorChar = " "
	r.Capitalize = CSFirst

	for i := 0; i < 20; i++ {
		p, err := r.Generate()
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
		sepRE := "\\Q" + r.SeparatorChar + "\\E"
		preCount := "{" + strconv.Itoa(r.Length-2) + "}"
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

func TestWLOneCap(t *testing.T) {
	threeG, err := NewWordList([]string{"once", "upon", "midnight", "dreary", "while", "pondered", "weak", "and", "weary", "over", "many"})
	if err != nil {
		t.Errorf("failed to create WL generator: %v", err)
	}
	// Test with random capitalization
	length := 5
	r := NewWLRecipe(length, threeG)
	r.SeparatorChar = " "
	r.Capitalize = CSOne

	tcWRE := "\\p{Lu}\\pL+"
	lcWRE := "\\p{Ll}\\pL+"
	wRE := "(?:" + tcWRE + ")|(?:" + lcWRE + ")"
	sepRE := "\\Q" + r.SeparatorChar + "\\E"
	preCount := "{" + strconv.Itoa(r.Length-1) + "}"
	leadRE := wRE + sepRE + "(?:" + wRE + sepRE + ")" + preCount
	res := "^" + leadRE + wRE + "$"
	re, err := regexp.Compile(res)

	if err != nil {
		t.Errorf("regexp %q did not compile: %v", res, err)
	}
	u, err := regexp.Compile("\\b\\p{Lu}")
	if err != nil {
		t.Errorf("regexp %q did not compile: %v", tcWRE, err)
	}
	l, err := regexp.Compile("\\b\\p{Ll}")
	if err != nil {
		t.Errorf("regexp %q did not compile: %v", lcWRE, err)
	}

	for i := 0; i < 10; i++ {
		p, err := r.Generate()
		ent := p.Entropy()
		expectedEnt := float32(19.619086) // 5 * log2(11) + log2(5)
		if err != nil {
			t.Errorf("failed to generate %d word password: %v", length, err)
		}
		if cmpFloat32(ent, expectedEnt, entCompTolerance) != 0 {
			t.Errorf("expected entropy (%.6f) != returned entropy (%.6f)", expectedEnt, ent)
		}

		pw := p.String()

		if !re.MatchString(pw) {
			t.Errorf("%q doesn't match %s", pw, re)
		}

		lCount := len(l.FindAllString(pw, -1)) // This appears to be really slow
		if lCount != r.Length-1 {
			t.Errorf("%d lowercase words in %q. Expected %d", lCount, pw, r.Length-1)
		}
		uCount := len(u.FindAllString(pw, -1))
		if uCount != 1 {
			t.Errorf("%d uppercase words in %q. Expected 1", uCount, pw)
		}
	}

}

func TestWLRandCapitalDistribution(t *testing.T) {

	if !doFallibleTests {
		t.Skipf("Skipping statistically fallible test: %v", t.Name())
	}

	threeG, err := NewWordList([]string{"egy", "kettő", "három"})
	if err != nil {
		t.Errorf("failed to create WL generator: %v", err)
	}
	length := 1024 // big enough to make misses unlikely, round enough for me to do math easily
	r := NewWLRecipe(length, threeG)
	r.SeparatorChar = " "
	r.Capitalize = CSRandom
	p, _ := r.Generate()
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

func TestNonLetterWL(t *testing.T) {
	cl := []string{"正確", "馬", "電池", "釘書針"}
	length := 5
	wl, err := NewWordList(cl)
	if err != nil {
		t.Errorf("failed to create wordlist generator from list %v: %v", wl, err)
	}
	a := NewWLRecipe(length, wl)
	a.SeparatorChar = " "
	a.Capitalize = CSOne

	// Because none of the words in the wordlist capitalize, the
	// a.Capitalize = CSOne setting makes no difference
	trueEnt := float32(math.Log2(float64(len(cl))) * float64(length))
	expectedEnt := trueEnt + float32(math.Log2(float64(length)))

	for i := 0; i < 20; i++ {
		p, err := a.Generate()
		pw, ent := p.String(), p.Entropy()
		if err != nil {
			t.Errorf("generator failed: %v", err)
		}

		// This test will fail if we use trueEnt instead of expected ent.
		// This is a consequence uppercasing some words making no difference
		if cmpFloat32(expectedEnt, ent, entCompTolerance) != 0 {
			t.Errorf("Expected entropy of %q is %.6f. Got %.6f", pw, expectedEnt, ent)
			t.Errorf("True entropy of %q is %.6f", pw, trueEnt)
		}
		// fmt.Println(pw)
	}
}

func TestSyllableDigit(t *testing.T) {
	// wl, err := NewWordList(abSyllables)
	wl, err := NewWordList([]string{"syl", "lab", "bull", "gen", "er", "at", "or"})
	if err != nil {
		t.Errorf("Couldn't create syllable generator: %v", err)
	}
	r := NewWLRecipe(12, wl)
	r.SeparatorFunc = SFDigits1
	r.Capitalize = CSOne

	// With a wordlist of 7 members, I get an expected entropy for these
	// attributes to be 48. int(12*log2(7) + log2(12) + 11*log2(10))
	expEnt := float32(73.81443)

	sylRE := "\\pL\\p{Ll}{1,3}" // A letter followed by 1-3 lowercase letters
	sepRE := "\\d"
	preCount := "{" + strconv.Itoa(r.Length-1) + "}"
	leadRE := "(?:" + sylRE + sepRE + ")" + preCount
	reStr := "^" + leadRE + sylRE + "$"
	re, err := regexp.Compile(reStr)
	if err != nil {
		t.Errorf("regexp %q did not compile: %v", re, err)
	}

	for i := 0; i < 20; i++ {
		p, err := r.Generate()
		pw, ent := p.String(), p.Entropy()
		if err != nil {
			t.Errorf("failed to generate syllable pw: %v", err)
		}
		// fmt.Println(pw)
		if !re.MatchString(pw) {
			t.Errorf("pwd %q didn't match regexp %q", pw, re)
		}
		if cmpFloat32(ent, expEnt, entCompTolerance) != 0 {
			t.Errorf("expected entropy of %.6f. Got %.6f", expEnt, ent)
		}
	}
}

func TestNonASCIISeparators(t *testing.T) {
	sl := []string{"uno", "dos", "tres"}
	length := 5
	wl, err := NewWordList(sl)
	if err != nil {
		t.Errorf("failed to create wordlist generator from list %v: %v", wl, err)
	}
	a := NewWLRecipe(length, wl)
	a.SeparatorChar = "¡"

	expectedEnt := float32(math.Log2(float64(len(sl))) * float64(length))

	for i := 0; i < 20; i++ {
		p, err := a.Generate()
		pw, ent := p.String(), p.Entropy()
		if err != nil {
			t.Errorf("generator failed: %v", err)
		}
		if cmpFloat32(expectedEnt, ent, entCompTolerance) != 0 {
			t.Errorf("Expected entropy of %q is %.6f. Got %.6f", pw, expectedEnt, ent)
		}
		// fmt.Println(pw)
	}
}
