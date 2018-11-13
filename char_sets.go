package spg

import (
	"fmt"
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

// requireFilter checks whether a candidate password has a character
// from each required character set
func requireFilter(pwd string, require reqSets) bool {
	if require == nil || len(require) == 0 {
		return true
	}
	for _, rset := range require {
		// ContainsAny does not treat empty strings like empty sets
		rs := stringFromSet(rset.s) // Separate var for debugging
		if rset.size() > 0 && !strings.ContainsAny(pwd, rs) {
			return false
		}
	}
	return true
}

func (r CharRecipe) fullAlphabet() (charList, error) {
	if r.allowedSet == nil {
		return nil, fmt.Errorf("allowedSet is nil")
	}
	fullABC := r.allowedSet.Union(r.requiredSets.union().s)
	return strings.Split(stringFromSet(fullABC), ""), nil
}

// disjointify trims all the required charsets
// and the allowed sets so that they are all mutually
// disjoint.
func disjointify(allowed charList, required reqSets) charList {
	abc := setFromString(strings.Join(allowed, ""))

	abc = abc.Difference(required.union().s)
	abcOut := strings.Split(stringFromSet(abc), "")

	return abcOut
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

/**
 ** Copyright 2018 AgileBits, Inc.
 ** Licensed under the Apache License, Version 2.0 (the "License").
 **/
