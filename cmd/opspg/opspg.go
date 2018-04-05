package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/agilebits/spg"
)

const (
	rtChar = iota
	rtWordlist
)

var flagRecipeType = flag.String("type", "char", "character (\"char\") or wordlist (\"wl\") recipe")
var flagLength = flag.Int("length", 16, "Length")

func main() {
	flag.Parse()

	recipeType := rtChar
	switch *flagRecipeType {
	case "wl":
		recipeType = rtWordlist
	default:
		recipeType = rtChar
	}
	switch recipeType {
	case rtWordlist:
		// We don't do this yet
		log.Fatal("word lists not yet impemented")

	case rtChar:
		doCharacterPassword()
	default:
		log.Fatalf("Unknown recipe type: %v\n", recipeType)
	}
}

func doCharacterPassword() {
	r := spg.CharRecipe{
		Length: *flagLength,
		Allow:  spg.Digits | spg.Letters | spg.Symbols,
	}

	pwd, err := r.Generate()
	if err != nil {
		log.Fatalf("Couldn't generate password: %v\n", err)
	}
	fmt.Println(pwd)
}
