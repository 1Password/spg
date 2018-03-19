package spg

import (
	"bytes"
	"testing"
)

func TestTokenizer(t *testing.T) {

	type tokenVec struct {
		Pwd        Password
		expectedTI TokenIndices
		expectedPW string
	}
	vecs := []tokenVec{}

	vecs = append(vecs, tokenVec{
		Pwd: Password{
			tokens: Tokens{
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
		expectedPW: "correct horse battery staple",
	})

	vecs = append(vecs, tokenVec{
		Pwd: Password{
			tokens: Tokens{
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
		expectedPW: "correct horse battery staple ",
	})

	vecs = append(vecs, tokenVec{
		Pwd: Password{
			tokens: Tokens{
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
		expectedPW: "P@ssw0rd1",
	})

	vecs = append(vecs, tokenVec{
		Pwd: Password{
			tokens: Tokens{
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
		expectedPW: "correcthorsebatterystaple",
	})

	for _, tVec := range vecs {

		tP := tVec.Pwd
		ts := tVec.Pwd.tokens
		ti, err := ts.TIndices()
		if err != nil {
			t.Errorf("failed to create token indices: %v", err)
		}
		if bytes.Compare(ti, tVec.expectedTI) != 0 {
			t.Errorf("ti != expected\n\tti: %v\n\tExpected: %v", ti, tVec.expectedTI)
		}
		pw := tP.String()
		ent := tP.Entropy
		if pw != tVec.expectedPW {
			t.Errorf("pw is %q. Expected %q", pw, tVec.expectedPW)
		}

		newP, err := Tokenize(pw, ti, ent)
		if err != nil {
			t.Errorf("couldn't tokenize: %v", err)
		}
		if newP.String() != pw {
			t.Errorf("%q should equal %q", newP.String(), pw)
		}
		if len(newP.tokens) != len(tP.tokens) {
			t.Errorf("tokens lengths don't match:\n\tOriginal: %v\n\tReconstructed: %v",
				tP.tokens, newP.tokens)
		} else { // only run this test if lengths are equal
			nt := newP.tokens
			for i, tok := range ts {
				if tok.Value() != nt[i].Value() {
					t.Errorf("%d-th tokens Values don't match: %q != %q", i, tok.Value(), nt[i].Value())
				}
				if tok.Type() != nt[i].Type() {
					t.Errorf("%d-th tokens Types don't match: %d != %d", i, tok.Type(), nt[i].Type())
				}
			}
		}
	}
}
