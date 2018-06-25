package spg

import (
	"strings"
	"testing"

	set "github.com/deckarep/golang-set"
)

func TestSetFromString(t *testing.T) {
	s := setFromString("abcd")
	if c := s.Cardinality(); c != 4 {
		t.Errorf("Set should be size 4, not %v", c)
	}

	for _, c := range []string{"a", "b", "c", "d"} {
		if !s.Contains(interface{}(c)) {
			t.Errorf("s should contain %q", c)
		}
	}
}

func TestStringFromSet(t *testing.T) {
	s := set.NewSet()

	for _, c := range []string{"a", "b", "c", "d"} {
		s.Add(interface{}(c))
	}
	str := stringFromSet(s)

	if len(str) != 4 {
		t.Errorf("string should be length 4, not %v", len(str))
	}

	for _, c := range []string{"a", "b", "c", "d"} {
		if !strings.Contains(str, c) {
			t.Errorf("str (%s) should contain %v", str, c)
		}
	}
}

func TestNewReqSet(t *testing.T) {
	s1 := newReqSet("abcabc", "TEST1")
	if s1 == nil {
		t.Error("Failed to create newReqSet")
	}

	if c := s1.s.Cardinality(); c != 3 {
		t.Errorf("Set should be size 3, not %v", c)
	}

	for _, c := range []string{"a", "b", "c"} {
		if !s1.s.Contains(interface{}(c)) {
			t.Errorf("s should contain %q", c)
		}
	}

	if s1.Name != "TEST1" {
		t.Errorf("Wrong name: %q", s1.Name)
	}
}

func TestReqSets(t *testing.T) {
	rs := make(reqSets, 2)
	rs[0] = *newReqSet("abcabc", "TEST1")
	rs[1] = *newReqSet("cde", "TEST2")

	union := rs.union()

	if size := union.s.Cardinality(); size != 5 {
		t.Errorf("Wrong union size: %d", size)
	}
}

func TestFilter(t *testing.T) {
	rs := make(reqSets, 2)
	rs[0] = *newReqSet("abcabc", "TEST1")
	rs[1] = *newReqSet("cde", "TEST2")

	type tvec struct {
		pwd      string
		expected bool
	}

	vecs := []tvec{
		{"xxxadxxxx", true},
		{"xxxxxx", false},
		{"xxxxaxxxx", false},
		{"xxxxxdxxx", false},
	}
	for _, v := range vecs {
		if res := includeFilter(v.pwd, rs); res != v.expected {
			t.Errorf("%q gets %v. Expected %v", v.pwd, res, v.expected)
		}

	}

}

func TestFilterEmpty(t *testing.T) {
	rs := make(reqSets, 3)
	rs[0] = *newReqSet("abcabc", "TEST1")
	rs[1] = *newReqSet("", "Empty")
	rs[2] = *newReqSet("cde", "TEST2")

	type tvec struct {
		pwd      string
		expected bool
	}

	vecs := []tvec{
		{"xxxadxxxx", true},
		{"xxxxxx", false},
		{"xxxxaxxxx", false},
		{"xxxxxdxxx", false},
	}
	for _, v := range vecs {
		if res := includeFilter(v.pwd, rs); res != v.expected {
			t.Errorf("%q gets %v. Expected %v", v.pwd, res, v.expected)
		}

	}

}

func TestBuildCharacterList(t *testing.T) {
	recip := &CharRecipe{Length: 10}
	recip.Allow = Digits | Symbols

	cl := recip.buildCharacterList()
	rs := recip.requeredSets

	t.Logf("len(cl): %d", len(cl))
	t.Logf("cl: %q", strings.Join(cl, ""))
	t.Logf("len(rs): %d", len(rs))

	for i := 0; i < len(rs); i++ {
		t.Logf("rs[%d].Name = %q", i, rs[i].Name)
		t.Logf("rs[%d].String() = %q", i, rs[i].String())
	}

}

func TestSetFromEmptyString(t *testing.T) {
	s := setFromString("")
	if c := s.Cardinality(); c != 0 {
		t.Errorf("Set should be size 0, not %v", c)
	}

}
