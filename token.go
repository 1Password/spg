package spg

import (
	"fmt"
	"math"
	"strings"
)

// TokenType holds the kinds of tokens within within a generated password
// Having it be an uint8 makes it easier to compact these into the token indices
type TokenType uint8

// For labelling tokens within a generated password
const (
	SeparatorType TokenType = iota
	AtomType
)

/*
Token is a unit within a generated password.
In "correct horse battery" there are five tokens
    []Token{
		{Value: "correct", Type: AtomType},
		{Value: " ", Type: SeparatorType},
		{Value: "horse", Type: AtomType},
		{Value: " ", Type: SeparatorType},
		{Value: "battery", Type: AtomType},
	}
*/
type Token struct {
	value string
	tType TokenType
}

// Value returns the string value
func (t Token) Value() string {
	return t.value
}

// Type returns the token type
func (t Token) Type() TokenType {
	return t.tType
}

// Tokens is the array of tokens that comprise a password
type Tokens []Token

// IndexKind is the kind of tokenization index.
// Token indices are compact byte arrays that can be
// used in conjunction to with a password string to reconstruct
// an array of Tokens
type IndexKind uint8

// Possible values for first byte of the compact token index index.
const (
	CharacterIndexKind   IndexKind = iota // Tokens are all Atoms and of length 1
	VarAtomsIndexKind                     // Tokes are all atoms (of potentally varying lengths)
	AlternatingIndexKind                  // Tokens are alternation of A S A S ... A
	FullIndexKind                         // Requires a full token index as sequeunce of token types is not predictable
)

func (ts Tokens) ofType(tType TokenType) []string {
	ret := []string{}
	for _, tok := range ts {
		if tok.tType == tType {
			ret = append(ret, tok.Value())
		}
	}
	return ret
}

// Indices can hold the indices needed to reconstruct tokens, separator from string
type Indices []byte

// MakeIndices returns a compact array of indices that indicate where a string is to be separated
// In the worst case it will need to encode both the length and the type of each token,
// thus requiring two bytes per token (plus the one leading byte)
// It does attempt to inspect the tokens to determine whether it can get away with
// encoding less information. The leading byte of the returned array contains necessary
// information about the particularly indexing used
//
// token lengths must be in (1, 255)
//
func (ts Tokens) MakeIndices() (Indices, error) {
	if len(ts) == 0 { // We aren't in a position to calculate this
		return nil, nil
	}

	kind := ts.Kind()
	switch kind {
	case CharacterIndexKind:
		return Indices{byte(kind)}, nil

	case AlternatingIndexKind:
		fallthrough
	case VarAtomsIndexKind:
		first := Indices{byte(kind)}
		ti := make(Indices, len(ts))
		for i, tok := range ts {
			v := tok.Value()
			lng := len(v)
			if lng > math.MaxUint8 {
				return nil, fmt.Errorf("token too large (%d)", lng)
			}
			ti[i] = uint8(lng)
		}
		// first + ti
		for _, t := range ti {
			first = append(first, t)
		}
		return first, nil

	default:
		first := Indices{byte(FullIndexKind)}
		ti := make(Indices, 2*len(ts))

		for i, tok := range ts {
			v := tok.Value()
			lng := len(v)
			tt := tok.Type()
			if lng > math.MaxUint8 {
				return nil, fmt.Errorf("token too large (%d)", lng)
			}
			ti[2*i] = uint8(lng)
			ti[(2*i)+1] = byte(tt)
		}

		// first + ti
		for _, t := range ti {
			first = append(first, t)
		}
		return first, nil
	}
}

// Kind looks at the tokens and works out what the most
// appropriate kind of token index we should use.
// It is exported in case callers would like help
// at guessing what kind of password they have,
// it does only provide a guess. It is no
// substitute for the original recipe
func (ts Tokens) Kind() IndexKind {

	// It's only atoms of length one (so character password)
	if ts.isAllAtoms() && ts.maxTokenLen() == 1 {
		return CharacterIndexKind
	}

	// Some atoms have length other than 1, so we will need
	// lengths in our index
	if ts.isAllAtoms() {
		return VarAtomsIndexKind
	}

	if ts.isAlternatingTokens() {
		return AlternatingIndexKind
	}

	// And when we don't know what other kind it is,
	return FullIndexKind
}

// isAlternatingTokes detects that kinds of passwords that were
// generated from a wordlist system with a non-empty separator
func (ts Tokens) isAlternatingTokens() bool {
	if len(ts)%2 != 1 {
		return false
	}
	types := ts.Types()
	if len(types) != 2 {
		return false
	}
	if !(types[AtomType] && types[SeparatorType]) {
		return false
	}
	for i, tok := range ts {
		tt := tok.Type()
		switch i % 2 {
		case 0: // evens should be Atoms
			if tt != AtomType {
				return false
			}
		case 1:
			if tt != SeparatorType {
				return false
			}
		}
	}
	return true
}

// Tokenize reconstructs a Password from a password string and Indices produced by MakeIndices()
func Tokenize(pw string, ti Indices, entropy float32) (Password, error) {
	p := Password{Entropy: entropy}
	chars := strings.Split(pw, "")

	if len(ti) == 0 {
		return p, fmt.Errorf("tokenization must begin with a TI Kind byte")
	}

	kind := IndexKind(ti[0])
	switch kind {
	case CharacterIndexKind:
		tokens := Tokens{}
		// all tokens are of type atom and are of length 1
		for _, c := range chars {
			tokens = append(tokens, Token{c, AtomType})
		}
		p.tokens = tokens
		return p, nil

	case VarAtomsIndexKind:
		tokens := make([]Token, len(ti)-1)
		prevPos := 0
		for i, tl := range ti[1:] {
			newPos := prevPos + int(tl)
			if newPos > len(chars) {
				return p, fmt.Errorf("password too short for indices")
			}
			v := strings.Join(chars[prevPos:newPos], "")
			tokens[i] = Token{v, AtomType}
			prevPos = newPos
		}
		p.tokens = tokens
		return p, nil

	case AlternatingIndexKind:
		tokens := make([]Token, len(ti)-1)
		prevPos := 0

		for i, tl := range ti[1:] {
			newPos := prevPos + int(tl)
			if newPos > len(chars) {
				return p, fmt.Errorf("password too short for indices")
			}
			v := strings.Join(chars[prevPos:newPos], "")
			tt := AtomType
			if i%2 == 1 {
				tt = SeparatorType
			}
			tokens[i] = Token{v, tt}
			prevPos = newPos
		}
		p.tokens = tokens
		return p, nil

	case FullIndexKind:
		tokens := make([]Token, len(ti)/2)

		prevPos := 0
		for i := 1; i < len(ti); i += 2 {
			tl := int(ti[i])
			tt := int(ti[i+1])
			newPos := prevPos + tl
			if newPos > len(chars) {
				return p, fmt.Errorf("password too short for indices")
			}
			v := strings.Join(chars[prevPos:newPos], "")
			tokens[i/2] = Token{v, TokenType(tt)}
			prevPos = newPos
		}
		p.tokens = tokens
		return p, nil
	default:
		return p, fmt.Errorf("Unknown TIIndex kind: %d", kind)
	}
}

// Types returns a set of all of the token types used within a password
// This is exported to help those who want to do fancy colorful display
// of passwords.
func (ts Tokens) Types() map[TokenType]bool {
	found := make(map[TokenType]bool)
	for _, tok := range ts {
		found[tok.Type()] = true
	}
	return found
}

// isAllOfType is true when all tokens in the password are of type tt
func (ts Tokens) isAllOfType(tt TokenType) bool {
	types := ts.Types()
	if len(types) == 1 && types[tt] {
		return true
	}
	return false
}

func (ts Tokens) maxTokenLen() int {
	max := 0
	for _, t := range ts {
		l := len(t.Value())
		if l > max {
			max = l
		}
	}
	return max
}

// isAllAtoms returns true when all of tokens are Atoms.
// It returns false if there are no tokens.
func (ts Tokens) isAllAtoms() bool { return ts.isAllOfType(AtomType) }

// Atoms returns the tokens (words, syllables, characters) that compromise the bulk of the password
func (ts Tokens) Atoms() []string { return ts.ofType(AtomType) }

// Separators are the separators between tokens.
// If this list is shorter than one less then the number of tokens, the last separator listed
// is used repeatedly to separate subsequent tokens.
// If this is nil, it is taken as nil no separators between tokens
func (ts Tokens) Separators() []string { return ts.ofType(SeparatorType) }

/**
 ** Copyright 2018 AgileBits, Inc.
 ** Licensed under the Apache License, Version 2.0 (the "License").
 **/
