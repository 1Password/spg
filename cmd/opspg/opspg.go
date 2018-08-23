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

// This really should use subcommands, as the flags for char and wl will
// have different meanings and defaults. But I haven't read up on how to do that.
var flagRecipeType = flag.String("type", "char", "character (\"char\") or wordlist (\"wl\") recipe")
var flagLength = flag.Int("length", 16, "Length")
var flagWLFile = flag.String("listfile", "", "Wordlist file")
var flagList = flag.String("list", "words", "Agilewords (\"words\") or AgileSyllables (\"syl\")")
var flagNumber = flag.Int("n", 1, "number of passwords to generate")
var flagEntropy = flag.Bool("e", false, "Display entropy")

func main() {
	flag.Parse()

	recipeType := rtChar
	switch *flagRecipeType {
	case "wl":
		recipeType = rtWordlist
	default:
		recipeType = rtChar
	}

	f := doWordlistPassword
	switch recipeType {
	case rtWordlist:
		f = doWordlistPassword

	case rtChar:
		f = doCharacterPassword
	default:
		log.Fatalf("Unknown recipe type: %v\n", recipeType)
	}
	for i := 0; i < *flagNumber; i++ {
		f()
	}
}

// really should use subcommands
func doCharacterPassword() {
	r := spg.CharRecipe{
		Length: *flagLength,
		Allow:  spg.Digits | spg.Letters | spg.Symbols,
	}

	pwd, err := r.Generate()
	if err != nil {
		log.Fatalf("Couldn't generate password: %v\n", err)
	}
	if *flagEntropy {
		fmt.Printf("%.2f:\t%v\n", pwd.Entropy, pwd)
	} else {
		fmt.Println(pwd)
	}
}

func doWordlistPassword() {
	if len(*flagWLFile) != 0 {
		log.Fatal("File reading not yet implemented")
	}

	var words []string
	switch *flagList {
	case "word":
		words = spg.AgileWords
	case "syl":
		words = spg.AgileSyllables
	default:
		log.Fatalf("list must be either %q or %q\n", "word", "syl")
	}

	wl, err := spg.NewWordList(words)
	if err != nil {
		log.Fatalf("Failed initiate wordlist: %v\n", err)
	}

	// Will need more command line options for settings, but lets just
	// hardcode stuff for now

	r := spg.NewWLRecipe(*flagLength, wl)
	r.SeparatorFunc = spg.SFDigits1
	r.Capitalize = spg.CSOne

	pwd, err := r.Generate()
	if err != nil {
		log.Fatalf("Couldn't generate password: %v\n", err)
	}

	if *flagEntropy {
		fmt.Printf("%.2f:\t%v\n", pwd.Entropy, pwd)
	} else {
		fmt.Println(pwd)
	}
}
