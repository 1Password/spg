package spg

import (
	"strings"

	set "github.com/deckarep/golang-set"
)

// reqSet is how we internally represent a set of characters
// which must be included
type reqSet struct {
	Name string  // To report which set wasn't matched on failure
	s    set.Set // this is the actual set.
}

type reqSets []reqSet

// incudeFilter checks whether a candidate password has a character
// from each required/include character set
func includeFilter(pwd string, include reqSets) bool {
	if include == nil || len(include) == 0 {
		return true
	}
	for _, r := range include {
		if !strings.ContainsAny(pwd, r.s.String()) {
			return false
		}
	}
	return true
}

// disjointify trims all the required charsets
// and the allowed sets so that they are all mutually
// disjoint.
func disjointify(allowed charList, include reqSets) (charList, reqSets) {
	abc := setFromString(strings.Join(allowed, ""))
	incOut := make(reqSets, 0)

	for i := 0; i < len(include)-1; i++ {
		abc = abc.Difference(include[i].s)
		for j := i + 1; j < len(include); j++ {
			incOut[j].s = include[j].s.Difference(include[i].s)
		}
	}

	abcOut := strings.Split(stringFromSet(abc), "")
	return abcOut, incOut
}

// setFromString creates a set of runes from a string
func setFromString(s string) set.Set {
	out := set.NewSet()

	for _, r := range strings.Split(s, "") {
		var i interface{} = r
		out.Add(i)
	}
	return out
}

func newReqSet(s, name string) *reqSet {
	r := &reqSet{
		Name: name,
		s:    setFromString(s),
	}
	return r
}

// stringFromSet will panic if the set isn't a set of strings
func stringFromSet(strSet set.Set) string {
	if strSet == nil {
		return ""
	}
	out := ""
	for r := range strSet.Iter() {
		out += r.(string)
	}
	return out
}

// String for reqSet will panic if the set's elements
// aren't all strings
func (r reqSet) String() string {
	if r.s == nil || r.s.Cardinality() == 0 {
		return ""
	}
	return stringFromSet(r.s)
}

func (rs reqSets) union() reqSet {
	u := reqSet{
		Name: "UNION",
		s:    set.NewSet(),
	}
	for _, s := range rs {
		u.s = u.s.Union(s.s)
	}
	return u
}

func (rs reqSets) size() int {
	if len(rs) == 0 {
		return 0
	}
	return rs.union().s.Cardinality()
}

func (r reqSet) size() int {
	if r.s == nil {
		return 0
	}
	return r.s.Cardinality()
}

func (rs reqSets) disjointify() reqSets {
	incOut := make(reqSets, 0)

	for i := 0; i < len(rs)-1; i++ {
		for j := i + 1; j < len(rs); j++ {
			incOut[j].s = rs[j].s.Difference(rs[i].s)
		}
	}
	return incOut
}
