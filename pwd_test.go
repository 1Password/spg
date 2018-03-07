package spg

import (
	"bufio"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

const doFallibleTests = false
const entCompTolerance = 10000 // entropy must be correct to 1 part in ten thousand

var abWords []string                          // this is where we will put the wordlist for testing
const wordsFilePath = "testdata/AgileWords.txt" // relative to where test in invoked

var abSyllables []string                             // this is where we will put the wordlist for testing
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

func TestSimpleEntropy(t *testing.T) {
	const Epsilon = 0.001 // margin for rounding errors in entropy calculation
	type ESVec struct {
		Length   int
		NElem    int
		Expected float64
	}

	vectors := []ESVec{}
	vectors = append(vectors, ESVec{1, 1024, 10.0})
	vectors = append(vectors, ESVec{5, 1024, 50.0})
	vectors = append(vectors, ESVec{5, 1, 0.0})
	vectors = append(vectors, ESVec{0, 1024, 0.0})
	vectors = append(vectors, ESVec{5, 18300, 70.79778014})
	vectors = append(vectors, ESVec{5, -1, math.NaN()})
	vectors = append(vectors, ESVec{5, 0, math.Inf(-1)})
	vectors = append(vectors, ESVec{-5, 1024, -50.0})
	vectors = append(vectors, ESVec{-5, 0, math.Inf(1)})

	for _, v := range vectors {
		if res := entropySimple(v.Length, v.NElem); cmpFloat(v.Expected, res, entCompTolerance) != 0 {
			t.Errorf("entropySimple(%d, %d) should be %.6f, not %.6f",
				v.Length, v.NElem, v.Expected, res)
		}
	}
}

// Now time for some character password tests

func TestDigitGenerator(t *testing.T) {
	g := new(CharacterPasswordGenerator)

	type ExIncVec struct {
		exc string
		inc string
		re  string
		ent float32
	}
	vectors := []ExIncVec{}

	vectors = append(vectors, ExIncVec{exc: "", inc: "", re: "^\\d{12}$", ent: 39.863137})
	vectors = append(vectors, ExIncVec{exc: "A8", inc: "", re: "^[0-79]{12}$", ent: 38.039100})
	vectors = append(vectors, ExIncVec{exc: "A8", inc: "08ABCDEF", re: "^[01-79B-F]{12}$", ent: 45.688259})

	for _, v := range vectors {
		re, err := regexp.Compile(v.re)
		if err != nil {
			t.Errorf("%q did not compile: %v", v.re, err)
		}
		attrs := NewGenAttrs(12)

		// Starting with digits-only
		attrs.ExcludeAmbiguous = false
		attrs.AllowDigit = true
		attrs.AllowLetter = false
		attrs.AllowSymbol = false

		attrs.ExcludeExtra = v.exc
		attrs.IncludeExtra = v.inc

		for i := 1; i <= 20; i++ {
			p, err := g.Generate(*attrs)
			pw, ent := p.String(), p.Entropy()
			if err != nil {
				t.Errorf("failed to generate password: %v", err)
			}
			if cmpFloat32(ent, v.ent, entCompTolerance) != 0 {
				t.Errorf("Expected entropy %.6f. Got %.6f instead", v.ent, ent)
			}
			if !re.MatchString(pw) {
				t.Errorf("%q didn't match %v", pw, re)
			}
		}
	}
}

func TestSyllableDigit(t *testing.T) {
	// g, err := NewWordListPasswordGenerator(abSyllables)
	g, err := NewWordListPasswordGenerator([]string{"syl", "lab", "bull", "gen", "er", "at", "or"})
	if err != nil {
		t.Errorf("Couldn't create syllable generator: %v", err)
	}
	attrs := NewGenAttrs(12)
	attrs.SeparatorFunc = SFDigits1
	attrs.Capitalize = CSOne

	// With a wordlist of 7 members, I get an expected entropy for these
	// attributes to be 48. int(12*log2(7) + log2(12) + 11*log2(10))
	expEnt := float32(73.81443)

	sylRE := "\\pL\\p{Ll}{1,3}" // A letter followed by 1-3 lowercase letters
	sepRE := "\\d"
	preCount := "{" + strconv.Itoa(attrs.Length-1) + "}"
	leadRE := "(?:" + sylRE + sepRE + ")" + preCount
	reStr := "^" + leadRE + sylRE + "$"
	re, err := regexp.Compile(reStr)
	if err != nil {
		t.Errorf("regexp %q did not compile: %v", re, err)
	}

	for i := 0; i < 20; i++ {
		p, err := g.Generate(*attrs)
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

func TestNonASCII(t *testing.T) {
	g := new(CharacterPasswordGenerator)
	length := 10
	a := NewGenAttrs(length)
	a.AllowDigit = false
	a.AllowLetter = false
	a.AllowSymbol = false
	a.IncludeExtra = "űβ™λ∞⊕💩"
	expectedEnt := float32(math.Log2(7.0) * float64(length))

	for i := 0; i < 20; i++ {
		p, err := g.Generate(*a)
		pw, ent := p.String(), p.Entropy()
		if err != nil {
			t.Errorf("Couldn't generate poopy password: %v", err)
		}
		// len(string) returns bytes not characters
		if gLength := len(strings.Split(pw, "")); gLength != length {
			t.Errorf("%q should be %d glyphs long, not %d", pw, length, gLength)
		}
		// fmt.Println(pw)
		if cmpFloat32(ent, expectedEnt, entCompTolerance) != 0 {
			t.Errorf("expected entropy of %.6f. Got %.6f", expectedEnt, ent)
		}
	}

}

func TestNonASCIISeparators(t *testing.T) {
	wl := []string{"uno", "dos", "tres"}
	length := 5
	g, err := NewWordListPasswordGenerator(wl)
	if err != nil {
		t.Errorf("failed to create wordlist generator from list %v: %v", wl, err)
	}
	a := NewGenAttrs(length)
	a.SeparatorChar = "¡"

	expectedEnt := float32(math.Log2(float64(len(wl))) * float64(length))

	for i := 0; i < 20; i++ {
		p, err := g.Generate(*a)
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

func TestNonLetterWL(t *testing.T) {
	wl := []string{"正確", "馬", "電池", "釘書針"}
	length := 5
	g, err := NewWordListPasswordGenerator(wl)
	if err != nil {
		t.Errorf("failed to create wordlist generator from list %v: %v", wl, err)
	}
	a := NewGenAttrs(length)
	a.SeparatorChar = " "
	a.Capitalize = CSOne

	// Because none of the words in the wordlist capitalize, the
	// a.Capitalize = CSOne setting makes no difference
	trueEnt := float32(math.Log2(float64(len(wl))) * float64(length))
	expectedEnt := trueEnt + float32(math.Log2(float64(length)))

	for i := 0; i < 20; i++ {
		p, err := g.Generate(*a)
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

// cmpFloat32 compares floats to 1 part in tolerance
// Returns 0 if equal (to 1 part in tolerance)
// 1 if a > b
// -1 if a < b
func cmpFloat32(a, b float32, tolerance int) int {

	// for some reason float32(int) doesn't exist. So we will do all of this
	// in float64

	return cmpFloat(float64(a), float64(b), tolerance)
}

func cmpFloat(a, b float64, tolerance int) int {
	tInv := 1.0 / math.Abs(float64(tolerance))
	avg := (math.Abs(a) + math.Abs(b)) / 2.0

	var epsilon float64
	if avg < tInv {
		epsilon = tInv
	} else {
		epsilon = avg * tInv
	}

	pInf := math.Inf(1)
	nInf := math.Inf(-1)

	if a == pInf && b == pInf {
		return 0
	}
	if a == nInf && b == nInf {
		return 0
	}
	// There is no good answer in this case, but for some of the
	// tests we run we do want to check NaN == NaN
	if math.IsNaN(a) && math.IsNaN(b) {
		return 0
	}
	if math.Abs(a-b) < epsilon {
		return 0
	}
	if a < b {
		return -1
	}
	return 1
}
