package main

import (
	"go.1password.io/spg"
)

type charRecipe struct {
	length  int
	allow   []string
	require []string
	exclude []string
}

var defaultCharRecipe = charRecipe{
	length:  20,
	require: []string{"uppercase", "lowercase", "digits"},
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
}
