package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/agilebits/spg"
)

const (
	rtChar = iota
	rtWordlist
)

type characterclass string

var ccMap = map[string]spg.CTFlag{
	"uppercase": spg.Uppers,
	"lowercase": spg.Lowers,
	"digits":    spg.Digits,
	"symbols":   spg.Symbols,
	"ambiguous": spg.Ambiguous,
}

var separatorMap = map[string]spg.SFFunction{
	"hyphen":     createSeparatorFunc("-"),
	"space":      createSeparatorFunc(" "),
	"comma":      createSeparatorFunc(","),
	"period":     createSeparatorFunc("."),
	"underscore": createSeparatorFunc("_"),
	"digit":      spg.SFDigits1,
	"none":       spg.SFNone,
}

var capitalizeMap = map[string]spg.CapScheme{
	"none":   spg.CSNone,
	"first":  spg.CSFirst,
	"all":    spg.CSAll,
	"random": spg.CSRandom,
	"one":    spg.CSOne,
}

var recipeCommand = flag.NewFlagSet("recipe", flag.ExitOnError)
var wordlistCommand = flag.NewFlagSet("wordlist", flag.ExitOnError)
var charactersCommand = flag.NewFlagSet("characters", flag.ExitOnError)

// Character flags
var flagLength = charactersCommand.Int("length", defaultCharRecipe.length, "generate a password <n> characters in length (default: 20)")
var flagAllow = charactersCommand.String("allow", "", "allow characters from <characterclasses> (default: all)")
var flagRequire = charactersCommand.String("require", "", "require at least one character from <characterclasses> (default: none)")
var flagExclude = charactersCommand.String("exclude", "", "exclude all characters from <characterclasses> regardless of other settings (default: ambiguous)")

// Wordlist flags
var flagSize = wordlistCommand.Int("size", 4, "generate a password with <n> elements (default: 4)")
var flagWordList = wordlistCommand.String("list", "words", "use built-in <wordlist> (default: words)")
var flagSeparator = wordlistCommand.String("separator", "hyphen", "separate components with <separatorclass> (default: hyphen)")
var flagCapitalize = wordlistCommand.String("capitalize", "none", "capitalize password according to <scheme> (default: none)")

// charactersCommand.String("")

// This really should use subcommands, as the flags for char and wl will
// have different meanings and defaults. But I haven't read up on how to do that.
// var flagRecipeType = flag.String("type", "char", "character (\"char\") or wordlist (\"wl\") recipe")
// var flagWLFile = flag.String("listfile", "", "Wordlist file")
// var flagList = flag.String("list", "words", "Agilewords (\"words\") or AgileSyllables (\"syl\")")
// var flagNumber = flag.Int("n", 1, "number of passwords to generate")
// var flagEntropy = flag.Bool("e", false, "Display entropy")

func main() {
	// flag.Parse()
	// if len(os.Args) == 1 {
	// 	return
	// }

	switch os.Args[1] {
	case "recipe":
		recipeCommand.Parse(os.Args[2:])

		// recipe := parseRecipe(*flagRecipe)

		// pwd, _ := recipe.Generate()
		// fmt.Println(pwd.String())

	case "characters":
		charactersCommand.Parse(os.Args[2:])
		generateCharacterPassword()
	case "wordlist":
		wordlistCommand.Parse(os.Args[2:])
		generateWordListPassword()
	}

	// recipeType := rtChar
	// switch *flagRecipeType {
	// case "wl":
	// 	recipeType = rtWordlist
	// default:
	// 	recipeType = rtChar
	// }

	// f := doWordlistPassword
	// switch recipeType {
	// case rtWordlist:
	// 	f = doWordlistPassword

	// case rtChar:
	// 	f = doCharacterPassword
	// default:
	// 	log.Fatalf("Unknown recipe type: %v\n", recipeType)
	// }
	// for i := 0; i < *flagNumber; i++ {
	// 	f()
	// }
}

// really should use subcommands
// func doCharacterPassword() {
// 	r := spg.CharRecipe{
// 		Length: *flagLength,
// 		Allow:  spg.Digits | spg.Letters | spg.Symbols,
// 	}

// 	pwd, err := r.Generate()
// 	if err != nil {
// 		log.Fatalf("Couldn't generate password: %v\n", err)
// 	}
// 	if *flagEntropy {
// 		fmt.Printf("%.2f:\t%v\n", pwd.Entropy, pwd)
// 	} else {
// 		fmt.Println(pwd)
// 	}
// }

// func doWordlistPassword() {
// 	if len(*flagWLFile) != 0 {
// 		log.Fatal("File reading not yet implemented")
// 	}

// 	var words []string
// 	switch *flagList {
// 	case "word":
// 		words = spg.AgileWords
// 	case "syl":
// 		words = spg.AgileSyllables
// 	default:
// 		log.Fatalf("list must be either %q or %q\n", "word", "syl")
// 	}

// 	wl, err := spg.NewWordList(words)
// 	if err != nil {
// 		log.Fatalf("Failed initiate wordlist: %v\n", err)
// 	}

// 	// Will need more command line options for settings, but lets just
// 	// hardcode stuff for now

// 	r := spg.NewWLRecipe(*flagLength, wl)
// 	r.SeparatorFunc = spg.SFDigits1
// 	r.Capitalize = spg.CSOne

// 	pwd, err := r.Generate()
// 	if err != nil {
// 		log.Fatalf("Couldn't generate password: %v\n", err)
// 	}

// 	if *flagEntropy {
// 		fmt.Printf("%.2f:\t%v\n", pwd.Entropy, pwd)
// 	} else {
// 		fmt.Println(pwd)
// 	}
// }

func createSeparatorFunc(value string) spg.SFFunction {
	return func() (string, spg.FloatE) {
		return value, 0
	}
}

func parseCharacterClasses(value string, defaults []string) spg.CTFlag {
	var ccFlags spg.CTFlag
	var classes []string

	if value != "" {
		classes = strings.Split(
			strings.Replace(value, " ", "", -1),
			",",
		)
	} else {
		classes = defaults
	}

	for _, c := range classes {
		if ccFlag, ok := ccMap[c]; ok {
			ccFlags |= ccFlag
		}
	}

	return ccFlags
}

func parseRecipe(value string) spg.Generator {
	recipe, ok := recipes[value]
	if !ok {
		os.Exit(1)
	}
	return recipe
}

func parseWordList(value string) (*spg.WordList, error) {
	var wl []string
	switch value {
	case "words":
		wl = spg.AgileWords
	case "syllables":
		wl = spg.AgileSyllables
	}

	return spg.NewWordList(wl)
}

func parseSeparator(value string) spg.SFFunction {
	return separatorMap[value]
}

func parseCapitalize(value string) spg.CapScheme {
	return capitalizeMap[value]
}

func generateWordListPassword() {
	wl, _ := parseWordList(*flagWordList)
	recipe := spg.NewWLRecipe(*flagSize, wl)
	recipe.SeparatorFunc = parseSeparator(*flagSeparator)
	recipe.Capitalize = parseCapitalize(*flagCapitalize)

	pwd, _ := recipe.Generate()
	fmt.Println(pwd.String())
}

func generateCharacterPassword() {
	recipe := spg.NewCharRecipe(*flagLength)
	recipe.Allow = parseCharacterClasses(*flagAllow, defaultCharRecipe.allow)
	recipe.Include = parseCharacterClasses(*flagRequire, defaultCharRecipe.require)
	recipe.Exclude = parseCharacterClasses(*flagExclude, defaultCharRecipe.exclude)

	pwd, _ := recipe.Generate()
	fmt.Println(pwd.String())
}
