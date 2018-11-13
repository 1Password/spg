package spg

import (
	"bytes"
	"testing"
)

func TestTokenizer(t *testing.T) {

	type tokenVec struct {
		Pwd        Password
		expectedTI Indices
		expectedPW string
	}
	vecs := []tokenVec{}

	vecs = append(vecs, tokenVec{
		Pwd: Password{
			tokens: Tokens{
				{"correct", AtomType},
				{" ", SeparatorType},
				{"horse", AtomType},
				{" ", SeparatorType},
				{"battery", AtomType},
				{" ", SeparatorType},
				{"staple", AtomType},
			},
			Entropy: 44.0,
		},
		expectedTI: Indices{
			byte(AlternatingIndexKind),
			7, 1, 5, 1, 7, 1, 6,
		},
		expectedPW: "correct horse battery staple",
	})

	vecs = append(vecs, tokenVec{
		Pwd: Password{
			tokens: Tokens{
				{"correct", AtomType},
				{" ", SeparatorType},
				{"horse", AtomType},
				{" ", SeparatorType},
				{"battery", AtomType},
				{" ", SeparatorType},
				{"staple", AtomType},
				{" ", SeparatorType},
			},
			Entropy: 44.0,
		},
		expectedTI: Indices{
			byte(FullIndexKind),
			7, byte(AtomType),
			1, byte(SeparatorType),
			5, byte(AtomType),
			1, byte(SeparatorType),
			7, byte(AtomType),
			1, byte(SeparatorType),
			6, byte(AtomType),
			1, byte(SeparatorType),
		},
		expectedPW: "correct horse battery staple ",
	})

	vecs = append(vecs, tokenVec{
		Pwd: Password{
			tokens: Tokens{
				{"P", AtomType},
				{"@", AtomType},
				{"s", AtomType},
				{"s", AtomType},
				{"w", AtomType},
				{"0", AtomType},
				{"r", AtomType},
				{"d", AtomType},
				{"1", AtomType},
			},
			Entropy: 14.0,
		},
		expectedTI: Indices{byte(CharacterIndexKind)},
		expectedPW: "P@ssw0rd1",
	})

	vecs = append(vecs, tokenVec{
		Pwd: Password{
			tokens: Tokens{
				{"correct", AtomType},
				{"horse", AtomType},
				{"battery", AtomType},
				{"staple", AtomType},
			},
			Entropy: 44.0,
		},
		expectedTI: Indices{
			byte(VarAtomsIndexKind),
			7, 5, 7, 6,
		},
		expectedPW: "correcthorsebatterystaple",
	})

	for _, tVec := range vecs {

		tP := tVec.Pwd
		ts := tVec.Pwd.Tokens()
		ti, err := ts.MakeIndices()
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
				tP.Tokens(), newP.Tokens())
		} else { // only run this test if lengths are equal
			nt := newP.Tokens()
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

/**
 ** Copyright 2018 AgileBits, Inc.
 ** Licensed under the Apache License, Version 2.0 (the "License").
 **/
