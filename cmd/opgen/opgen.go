package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"go.1password.io/spg"
)

const (
	rtChar = iota
	rtWordlist
)

// Exit statuses for os.Exit(). Follow narrow Unix convenstions (1-127 for errors, 0 for success)
const (
	ExitSuccess  = iota // Success must be 0
	ExitCatchall        // Catch all should be 1. For all otherwise unspecified errors
	ExitUsage           // Usage errors.
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

// Subcommands
var recipeCommand = flag.NewFlagSet("recipe", flag.ExitOnError)
var wordlistCommand = flag.NewFlagSet("wordlist", flag.ExitOnError)
var charactersCommand = flag.NewFlagSet("characters", flag.ExitOnError)

// Character flags
var flagLength = charactersCommand.Int("length", defaultCharRecipe.length, "generate a password <n> characters in length")
var flagAllow = charactersCommand.String("allow", "", "allow characters from <characterclasses>")
var flagRequire = charactersCommand.String("require", "", "require at least one character from <characterclasses>")
var flagExclude = charactersCommand.String("exclude", "", "exclude all characters from <characterclasses> regardless of other settings")
var flagEntropyCR = charactersCommand.Bool("entropy", false, "show the entropy of the password")
var flagNumberCR = charactersCommand.Int("number", 1, "generate <n> passwords with the same recipe")

// Wordlist flags
var flagSize = wordlistCommand.Int("size", 4, "generate a password with <n> elements")
var flagWordList = wordlistCommand.String("list", "words", "use built-in <wordlist>")
var flagWordListFile = wordlistCommand.String("file", "", "use a wordlist file at the specified <path>")
var flagSeparator = wordlistCommand.String("separator", "hyphen", "separate components with <separatorclass>")
var flagCapitalize = wordlistCommand.String("capitalize", "none", "capitalize password according to <scheme>")
var flagEntropyWL = wordlistCommand.Bool("entropy", false, "show the entropy of the password")
var flagNumberWL = wordlistCommand.Int("number", 1, "generate <n> passwords with the same recipe")

func main() {
	flag.Parse()
	if len(os.Args) == 1 {
		printUsage()
		os.Exit(ExitUsage)
	}

	var generator spg.Generator
	switch os.Args[1] {
	case "recipe":
		recipeCommand.Parse(os.Args[2:])

		// recipe := parseRecipe(*flagRecipe)

		// pwd, _ := recipe.Generate()
		// fmt.Println(pwd.String())

	case "characters":
		charactersCommand.Parse(os.Args[2:])
		generator = charGenerator()
	case "wordlist":
		wordlistCommand.Parse(os.Args[2:])
		generator = wlGenerator()
	default:
		printUsage()
		os.Exit(ExitUsage)
	}

	for i := 1; i < *flagNumberCR+*flagNumberWL; i++ {
		pwd, err := generator.Generate()
		if err != nil {
			log.Fatalln("Error generating password:", err)
			return
		}

		if *flagEntropyWL || *flagEntropyCR {
			fmt.Printf("%.2f:\t%s\n", generator.Entropy(), pwd.String())
		} else {
			fmt.Println(pwd.String())
		}
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
		os.Exit(ExitUsage)
	}
	return recipe
}

func parseWordList(value string) *spg.WordList {
	var words []string
	switch value {
	case "words":
		words = spg.AgileWords
	case "syllables":
		words = spg.AgileSyllables
	default:
		printUsage()
		os.Exit(ExitUsage)
	}

	wordList, _ := spg.NewWordList(words)
	return wordList
}

func parseSeparator(value string) spg.SFFunction {
	return separatorMap[value]
}

func parseCapitalize(value string) spg.CapScheme {
	return capitalizeMap[value]
}

func wlGenerator() *spg.WLRecipe {
	var wl *spg.WordList
	if *flagWordListFile != "" {
		wl = loadWordListFile(*flagWordListFile)
	} else {
		wl = parseWordList(*flagWordList)
	}

	recipe := spg.NewWLRecipe(*flagSize, wl)
	recipe.SeparatorFunc = parseSeparator(*flagSeparator)
	recipe.Capitalize = parseCapitalize(*flagCapitalize)

	return recipe
}

func charGenerator() *spg.CharRecipe {
	recipe := spg.NewCharRecipe(*flagLength)
	recipe.Allow = parseCharacterClasses(*flagAllow, defaultCharRecipe.allow)
	recipe.Require = parseCharacterClasses(*flagRequire, defaultCharRecipe.require)
	recipe.Exclude = parseCharacterClasses(*flagExclude, defaultCharRecipe.exclude)

	return recipe
}

func loadWordListFile(path string) *spg.WordList {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln("Error opening file:", path, err)
	}

	words := strings.Fields(string(data))
	wordList, err := spg.NewWordList(words)
	if err != nil {
		log.Fatalln("Error creating wordlist:", err)
	}
	return wordList
}

func printUsage() {
	fmt.Println(`
opgen recipe [<recipe> | --file=<recipefile>] [--number=<n>] [--entropy]

	--file      use a recipe file at the specified path
	--number    generate <n> passwords with the same recipe (default: 1)
	--entropy   show the entropy of the password

	<recipe>: memorable, syllables, pin

opgen characters [--length=<n>] [--allow=<characterclasses>]
				[--exclude=<characterclasses>] [--require=<characterclasses>]
				[--number=<numberofpasswords>] [--entropy]

	--length    generate a password <n> characters in length (default: 20)
	--allow     allow characters from <characterclasses> (default: all)
	--exclude   exclude all characters from <characterclasses> regardless of
					other settings (default: ambiguous)
	--require   require at least one character from <characterclasses>
					(default: none)
	--number    generate <n> passwords with the same recipe (default: 1)
	--entropy   show the entropy of the password

	<characterclasses>: uppercase, lowercase, digits, symbols, ambiguous

opgen wordlist [--list=<wordlist> | --file=<wordlistfile>] [--size=<n>]
				[--separator=<separatorclass>] [--capitalize=<scheme>]
				[--number=<numberofpasswords>] [--entropy]

	--list         use built-in <wordlist> (default: words)
	--file         use a wordlist file at the specified path
	--size         generate a password with <n> elements (default: 4)
	--separator    separate components with <separatorclass> (default: hyphen)
	--capitalize   capitalize password according to <scheme> (default: none)
	--number       generate <n> passwords with the same recipe (default: 1)
	--entropy      show the entropy of the password

	<wordlist>: words, syllables
	<separatorclass>: hyphen, space, comma, period, underscore, digit, none
	capitalization <scheme>: none, first, all, random, one
	`)
}
