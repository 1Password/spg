package main

import (
	"github.com/1password/spg"
)

type charRecipe struct {
	length  int
	allow   []string
	require []string
	exclude []string
}

var defaultCharRecipe = charRecipe{
	length:  20,
	allow:   []string{"uppercase", "lowercase", "digits", "symbols"},
	exclude: []string{"ambiguous"},
}

var recipes = map[string]spg.Generator{
	"pin": &spg.CharRecipe{
		Length: 4,
		Allow:  spg.Digits,
	},
	"memorable": &spg.WLRecipe{
		Length:        4,
		SeparatorChar: "-",
	},
	"syllables": &spg.WLRecipe{
		Length: 5,
	},
	// compromise attempts to fix strength, memorability, and type-ability
	"compromise": &spg.WLRecipe{
		Length:        4,
		Capitalize:    spg.CSOne,
		SeparatorFunc: spg.SFDigits1,
	},
}
