package spg

import (
	rand "crypto/rand"
	"encoding/binary"
	"fmt"
	"math"
	"strings"
)

// subtractString returns a copy of source with any characters that appear in remove removed.
// It does not presever order.
func subtractString(source, remove string) string {

	src := setFromString(source)
	rm := setFromString(remove)

	diff := src.Difference(rm)
	out := stringFromSet(diff)
	return out
}

// nFromString picks characters from a sting. This is for internal use only. It does not check for duplicates in the string
func nFromString(ab string, n int) (string, float64) {
	if len(ab) == 0 {
		return "", 0.0
	}
	if n < 1 {
		return "", 0.0
	}
	ent := float64(n) * math.Log2(float64(len(ab)))
	sep := ""
	rAB := strings.Split(ab, "") // an AlphaBet of runes
	for i := 1; i <= n; i++ {
		sep += string(rAB[randomUint32n(uint32(len(rAB)))])
	}
	return sep, ent

}

// randomUint32 creates a random 32 bit unsigned integer
func randomUint32() uint32 {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		panic("PRNG gen error:" + err.Error())
	}

	return binary.BigEndian.Uint32(b)
}

// entropySimple takes the password length and the number of elements in the alphabet
// (nelem would be number of words in a wordlist or number of characters in the alphabet
// from which the password is generated).
// It returns a float64
func entropySimple(length int, nelem int) FloatE {
	// The entropy of, say, a length character password
	// with characters drawn for letters and digits (so nelem is 62) would be
	// length * log2(62).

	if nelem < 1 {
		// We will end up returning NaN or -Inf, so we are only logging here
		fmt.Printf("entropySimple: There must be a positive number of elements. Not %d\n", nelem)
	}
	entPerUnit := math.Log2(float64(nelem))
	return FloatE(float64(length) * entPerUnit)
}

// randomUint32n returns, as a uint32, a non-negative random number in [0,n) from a cryptographic appropriate source.
// It panics if a security-sensitive random number cannot be created.
// Care is taken to avoid modulo bias.
//
// Based on Int31n from the math/rand package..
func randomUint32n(n uint32) uint32 {
	if n <= 1 {
		return 0
	}
	if n&(n-1) == 0 { // n is power of two, can mask
		return randomUint32() & (n - 1)
	}
	discard := uint32(math.MaxUint32 - math.MaxUint32%n)
	v := randomUint32()
	for v >= discard {
		v = randomUint32()
	}
	return v % n
}

/**
 ** Copyright 2018 AgileBits, Inc.
 ** Licensed under the Apache License, Version 2.0 (the "License").
 **/
