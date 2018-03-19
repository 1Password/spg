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
	SeparatorTokenType TokenType = iota
	AtomTokenType
)

/*
Token is a unit within a generated password.
In "correct horse battery" there are five tokens
    []Token{
		{Value: "correct", Type: AtomTokenType},
		{Value: " ", Type: SeparatorTokenType},
		{Value: "horse", Type: AtomTokenType},
		{Value: " ", Type: SeparatorTokenType},
		{Value: "battery", Type: AtomTokenType},
	}
*/
type Token struct {
	Value string
	Type  TokenType
}

// TIndexKind is the kind of tokenization index.
// Token indices are compact byte arrays that can be
// used in conjunction to with a password string to reconstruct
// an array of Tokens
type TIndexKind uint8

// Possible values for first byte of the compact token index index.
const (
	CharacterTIIndexKind   TIndexKind = iota // Tokens are all Atoms and of length 1
	VarAtomsTIIndexKind                      // Tokes are all atoms (of potentally varying lengths)
	AlternatingTIIndexKind                   // Tokens are alternation of A S A S ... A
	FullTIIndexKind                          // Requires a full token index as sequeunce of token types is not predictable
)

func (p Password) tokensOfType(tType TokenType) []string {
	ret := []string{}
	for _, tok := range p.Tokens {
		if tok.Type == tType {
			ret = append(ret, tok.Value)
		}
	}
	return ret
}


// TokenIndices can hold the indices needed to reconstruct tokens, separator from string
type TokenIndices []byte

// TIndices returns a compact array of indices that indicate where a string is to be separated
// In the worst case it will need to encode both the length and the type of each token,
// thus requiring two bytes per token (plus the one leading byte)
// It does attempt to inspect the tokens to determine whether it can get away with
// encoding less information. The leading byte of the returned array contains necessary
// information about the particularly indexing used
//
// token lengths must be in (1, 255)
//
func (p Password) TIndices() (TokenIndices, error) {
	if len(p.Tokens) == 0 { // We aren't in a position to calculate this
		return nil, nil
	}

	kind := p.tokensKind()
	switch kind {
	case CharacterTIIndexKind:
		return TokenIndices{byte(kind)}, nil

	case AlternatingTIIndexKind:
		fallthrough
	case VarAtomsTIIndexKind:
		first := TokenIndices{byte(kind)}
		ti := make(TokenIndices, len(p.Tokens))
		for i, tok := range p.Tokens {
			v := tok.Value
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
		first := TokenIndices{byte(FullTIIndexKind)}
		ti := make(TokenIndices, 2*len(p.Tokens))

		for i, tok := range p.Tokens {
			v := tok.Value
			lng := len(v)
			tt := tok.Type
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

func (p Password) tokensKind() TIndexKind {

	// It's only atoms of length one (so character password)
	if p.IsAllAtoms() && p.maxTokenLen() == 1 {
		return CharacterTIIndexKind
	}

	// Some atoms have length other than 1, so we will need
	// lengths in our index
	if p.IsAllAtoms() {
		return VarAtomsTIIndexKind
	}

	if p.isAlternatingTokens() {
		return AlternatingTIIndexKind
	}

	// And when we don't know what other kind it is,
	return FullTIIndexKind
}

func (p Password) isAlternatingTokens() bool {
	toks := p.Tokens
	if len(toks)%2 != 1 {
		return false
	}
	types := p.TokenTypes()
	if len(types) != 2 {
		return false
	}
	if !(types[AtomTokenType] && types[SeparatorTokenType]) {
		return false
	}
	for i, tok := range toks {
		switch i % 2 {
		case 0: // evens should be Atoms
			if tok.Type != AtomTokenType {
				return false
			}
		case 1:
			if tok.Type != SeparatorTokenType {
				return false
			}
		}
	}
	return true
}

// Tokenize reconstructs a Password from a password string and TokenIndices produced by TIndices()
func Tokenize(pw string, ti TokenIndices, entropy float32) (Password, error) {
	p := Password{Entropy: entropy}
	chars := strings.Split(pw, "")

	if len(ti) == 0 {
		return p, fmt.Errorf("tokenization must begin with a TI Kind byte")
	}

	kind := TIndexKind(ti[0])
	switch kind {
	case CharacterTIIndexKind:
		toks := []Token{}
		// all tokens are of type atom and are of length 1
		for _, c := range chars {
			toks = append(toks, Token{c, AtomTokenType})
		}
		p.Tokens = toks
		return p, nil

	case VarAtomsTIIndexKind:
		toks := make([]Token, len(ti)-1)
		prevPos := 0
		for i, tl := range ti[1:] {
			newPos := prevPos + int(tl)
			if newPos > len(chars) {
				return p, fmt.Errorf("password too short for indices")
			}
			v := strings.Join(chars[prevPos:newPos], "")
			toks[i] = Token{v, AtomTokenType}
			prevPos = newPos
		}
		p.Tokens = toks
		return p, nil

	case AlternatingTIIndexKind:
		toks := make([]Token, len(ti)-1)
		prevPos := 0

		for i, tl := range ti[1:] {
			newPos := prevPos + int(tl)
			if newPos > len(chars) {
				return p, fmt.Errorf("password too short for indices")
			}
			v := strings.Join(chars[prevPos:newPos], "")
			tt := AtomTokenType
			if i%2 == 1 {
				tt = SeparatorTokenType
			}
			toks[i] = Token{v, tt}
			prevPos = newPos
		}
		p.Tokens = toks
		return p, nil

	case FullTIIndexKind:
		toks := make([]Token, len(ti)/2)

		prevPos := 0
		for i := 1; i < len(ti); i += 2 {
			tl := int(ti[i])
			tt := int(ti[i+1])
			newPos := prevPos + tl
			if newPos > len(chars) {
				return p, fmt.Errorf("password too short for indices")
			}
			v := strings.Join(chars[prevPos:newPos], "")
			toks[i/2] = Token{v, TokenType(tt)}
			prevPos = newPos
		}
		p.Tokens = toks
		return p, nil
	default:
		return p, fmt.Errorf("Unknown TIIndex kind: %d", kind)
	}
}

// TokenTypes returns a set of all of the token types used within a password
func (p Password) TokenTypes() map[TokenType]bool {
	found := make(map[TokenType]bool)
	for _, tok := range p.Tokens {
		found[tok.Type] = true
	}
	return found
}

// isAllOfType is true when all tokens in the password are of type tt
func (p Password) isAllOfType(tt TokenType) bool {
	types := p.TokenTypes()
	if len(types) == 1 && types[tt] {
		return true
	}
	return false
}

func (p Password) maxTokenLen() int {
	max := 0
	for _, t := range p.Tokens {
		if len(t.Value) > max {
			max = len(t.Value)
		}
	}
	return max
}

// IsAllAtoms returns true when all of tokens are Atoms.
// It returns false if there are no tokens
func (p Password) IsAllAtoms() bool { return p.isAllOfType(AtomTokenType) }
