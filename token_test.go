package spg

import (
	"bytes"
	"testing"
)

func TestTokenizer(t *testing.T) {

	type tokenVec struct {
		Pwd        Password
		expectedTI TokenIndices
	}
	vecs := []tokenVec{}

	vecs = append(vecs, tokenVec{
		Pwd: Password{
			Tokens: []Token{
				{"correct", AtomTokenType},
				{" ", SeparatorTokenType},
				{"horse", AtomTokenType},
				{" ", SeparatorTokenType},
				{"battery", AtomTokenType},
				{" ", SeparatorTokenType},
				{"staple", AtomTokenType},
			},
			Entropy: 44.0,
		},
		expectedTI: TokenIndices{
			byte(AlternatingTIIndexKind),
			7, 1, 5, 1, 7, 1, 6,
		},
	})

	vecs = append(vecs, tokenVec{
		Pwd: Password{
			Tokens: []Token{
				{"correct", AtomTokenType},
				{" ", SeparatorTokenType},
				{"horse", AtomTokenType},
				{" ", SeparatorTokenType},
				{"battery", AtomTokenType},
				{" ", SeparatorTokenType},
				{"staple", AtomTokenType},
				{" ", SeparatorTokenType},
			},
			Entropy: 44.0,
		},
		expectedTI: TokenIndices{
			byte(FullTIIndexKind),
			7, byte(AtomTokenType),
			1, byte(SeparatorTokenType),
			5, byte(AtomTokenType),
			1, byte(SeparatorTokenType),
			7, byte(AtomTokenType),
			1, byte(SeparatorTokenType),
			6, byte(AtomTokenType),
			1, byte(SeparatorTokenType),
		},
	})

	vecs = append(vecs, tokenVec{
		Pwd: Password{
			Tokens: []Token{
				{"P", AtomTokenType},
				{"@", AtomTokenType},
				{"s", AtomTokenType},
				{"s", AtomTokenType},
				{"w", AtomTokenType},
				{"0", AtomTokenType},
				{"r", AtomTokenType},
				{"d", AtomTokenType},
				{"1", AtomTokenType},
			},
			Entropy: 14.0,
		},
		expectedTI: TokenIndices{byte(CharacterTIIndexKind)},
	})

	vecs = append(vecs, tokenVec{
		Pwd: Password{
			Tokens: []Token{
				{"correct", AtomTokenType},
				{"horse", AtomTokenType},
				{"battery", AtomTokenType},
				{"staple", AtomTokenType},
			},
			Entropy: 44.0,
		},
		expectedTI: TokenIndices{
			byte(VarAtomsTIIndexKind),
			7, 5, 7, 6,
		},
	})

	for _, tVec := range vecs {

		tP := tVec.Pwd
		ti, err := tP.TIndices()
		if err != nil {
			t.Errorf("failed to create token indices: %v", err)
		}
		if bytes.Compare(ti, tVec.expectedTI) != 0 {
			t.Errorf("ti != expected\n\tti: %v\n\tExpected: %v", ti, tVec.expectedTI)
		}
		pw := tP.String()
		ent := tP.Entropy

		newP, err := Tokenize(pw, ti, ent)
		if err != nil {
			t.Errorf("couldn't tokenize: %v", err)
		}
		if newP.String() != pw {
			t.Errorf("%q should equal %q", newP.String(), pw)
		}
		if len(newP.Tokens) != len(tP.Tokens) {
			t.Errorf("tokens lengths don't match:\n\tOriginal: %v\n\tReconstructed: %v",
				tP.Tokens, newP.Tokens)
		} else { // only run this test if lengths are equal
			nt := newP.Tokens
			for i, tok := range tP.Tokens {
				if tok.Value != nt[i].Value {
					t.Errorf("%d-th tokens Values don't match: %q != %q", i, tok.Value, nt[i].Value)
				}
				if tok.Type != nt[i].Type {
					t.Errorf("%d-th tokens Types don't match: %d != %d", i, tok.Type, nt[i].Type)
				}
			}
		}
	}
}
