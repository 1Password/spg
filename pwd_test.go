package spg

import (
	"fmt"
	"math"
	"regexp"
	"strings"
	"testing"
)

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
		res := entropySimple(v.Length, v.NElem)
		if cmpFloat(v.Expected, float64(res), entCompTolerance) != 0 {
			t.Errorf("entropySimple(%d, %d) should be %.6f, not %.6f",
				v.Length, v.NElem, v.Expected, res)
		}
	}
}

// Now time for some character password tests

func TestDigitGenerator(t *testing.T) {

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
		r := &CharRecipe{Length: 12}

		// Starting with digits-only
		r.Allow = 0 | Digits
		r.Exclude = 0 // Don't exclude ambiguous

		r.ExcludeChars = v.exc
		r.AllowChars = v.inc

		for i := 1; i <= 20; i++ {
			p, err := r.Generate()
			pw, ent := p.String(), p.Entropy
			if err != nil {
				t.Errorf("failed to generate password: %v", err)
			}
			if cmpFloat32(ent, v.ent, entCompTolerance) != 0 {
				t.Errorf("Expected entropy %.6f. Got %.6f instead for %q", v.ent, ent, pw)
			}
			if !re.MatchString(pw) {
				t.Errorf("%q didn't match %v", pw, re)
			}
		}
	}
}

func TestNonASCII(t *testing.T) {
	length := 10
	r := &CharRecipe{
		Length:     length,
		Allow:      0,
		AllowChars: "űβ™λ∞⊕💩",
	}
	expectedEnt := float32(math.Log2(7.0) * float64(length))

	for i := 0; i < 20; i++ {
		p, err := r.Generate()
		pw, ent := p.String(), p.Entropy
		if err != nil {
			t.Errorf("Couldn't generate poopy password: %v", err)
		}
		// len(string) returns bytes not characters
		if gLength := len(strings.Split(pw, "")); gLength != length {
			t.Errorf("%q should be %d glyphs long, not %d", pw, length, gLength)
		}
		// fmt.Println(pw)
		if cmpFloat32(ent, expectedEnt, entCompTolerance) != 0 {
			t.Errorf("expected entropy of %.6f. Got %.6f for %q", expectedEnt, ent, pw)
		}
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

// Examples
// Because these have random output, we can't use "Output:",
func ExampleCharRecipe_pin() {
	r := CharRecipe{
		Length: 4,      // Password will be 4 characters long
		Allow:  Digits, // and comprised of digits
	}
	pwd, _ := r.Generate() // In real code, you would check error
	fmt.Println(pwd)
}

func ExampleCharRecipe_default() {
	r := NewCharRecipe(15) // 15 character passwords

	// Let's generate 5 passwords to get a small taste of them
	for i := 0; i < 5; i++ {
		p, _ := r.Generate() // You'd check for errors in real code
		fmt.Printf("Password: %q\tEntropy: %.3f\n", p, p.Entropy)
	}
}

func ExampleCharRecipe_lowerdigits() {
	r := CharRecipe{
		Length:  17,              // Password will be 17 characters long
		Allow:   Lowers | Digits, // and may contain lowercase letters and digits
		Exclude: Ambiguous,       // but no ambiguous characters
	}

	// Let's generate five of them for a small sample
	for i := 0; i < 5; i++ {
		p, _ := r.Generate() // You would check error in real code
		fmt.Printf("Password: %q\tEntropy: %.3f\n", p, p.Entropy)
	}
}

// This will run the CharRecipe examples if the -v flat is passed
// to go test
func TestExampleCharRecipe(t *testing.T) {
	if !testing.Verbose() {
		return
	}
	ExampleCharRecipe_default()
	ExampleCharRecipe_pin()
	ExampleCharRecipe_lowerdigits()
}

// Now let's play with inclusion requirements. Sadly, these may be statistical tests

func TestDigitInclusion(t *testing.T) {
	r := CharRecipe{
		Length:  6,
		Allow:   Letters,
		Include: Digits | Symbols,
	}

	successes := 0
	trials := 50
	for i := 0; i < trials; i++ {
		p, err := r.Generate()
		if err == nil {
			successes++
			if !strings.ContainsAny(p.String(), "0123456789") {
				t.Errorf("%q does not contain a digit", p.String())
			}
			if !strings.ContainsAny(p.String(), ctSymbols) {
				t.Errorf("%q does not contain a symbol", p.String())
			}
		}
	}
	if testing.Verbose() {
		fmt.Printf("Able to generate %d length strings %d/%d times\n", r.Length, successes, trials)
	}
}
