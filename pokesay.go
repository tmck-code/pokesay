package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/mitchellh/go-wordwrap"
	"github.com/tmck-code/pokesay-go/src/pokedex"
	"github.com/tmck-code/pokesay-go/src/timer"
)

var (
	//go:embed build/pokedex.gob
	GOBCategory []byte
	//go:embed build/*cow
	GOBCowData embed.FS
	Rand       rand.Source = rand.NewSource(time.Now().UnixNano())
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func randomInt(n int) int {
	return rand.New(Rand).Intn(n)
}

func printSpeechBubbleLine(line string, width int) {
	if len(line) > width {
		fmt.Println("|", line)
	} else {
		fmt.Println("|", line, strings.Repeat(" ", width-len(line)), "|")
	}
}

func printWrappedText(line string, width int, tabSpaces string) {
	for _, wline := range strings.Split(wordwrap.WrapString(strings.Replace(line, "\t", tabSpaces, -1), uint(width)), "\n") {
		printSpeechBubbleLine(wline, width)
	}
}

func printSpeechBubble(scanner *bufio.Scanner, args Args) {
	border := strings.Repeat("-", args.Width+2)
	fmt.Println("/" + border + "\\")

	for scanner.Scan() {
		line := scanner.Text()

		if !args.NoTabSpaces {
			line = strings.Replace(line, "\t", args.TabSpaces, -1)
		}
		if args.NoWrap {
			printSpeechBubbleLine(line, args.Width)
		} else {
			printWrappedText(line, args.Width, args.TabSpaces)
		}
	}
	fmt.Println("\\" + border + "/")
	for i := 0; i < 4; i++ {
		fmt.Println(strings.Repeat(" ", i+8), "\\")
	}
}

func printPokemon(choice *pokedex.PokemonEntry) {
	d, _ := GOBCowData.ReadFile(pokedex.EntryFpath(choice.Index))
	fmt.Printf("%s\nchoice: %s\n", pokedex.Decompress(d), choice.Name)
}

func chooseRandomCategory(keys [][]string, categories pokedex.PokemonTrie) []*pokedex.PokemonEntry {
	category, err := categories.GetCategory(keys[randomInt(len(keys)-1)])
	check(err)
	return category
}

func chooseRandomPokemon(pokemon []*pokedex.PokemonEntry) *pokedex.PokemonEntry {
	return pokemon[randomInt(len(pokemon))]
}

type Args struct {
	Width          int
	NoWrap         bool
	TabSpaces      string
	NoTabSpaces    bool
	ListCategories bool
	Category       string
	NameToken      string
}

func parseFlags() Args {
	width := flag.Int("width", 80, "the max speech bubble width")
	noWrap := flag.Bool("nowrap", false, "disable text wrapping (fastest)")
	tabWidth := flag.Int("tabwidth", 4, "replace any tab characters with N spaces")
	noTabSpaces := flag.Bool("notabspaces", false, "do not replace tab characters (fastest)")
	fastest := flag.Bool("fastest", false, "run with the fastest possible configuration (-nowrap -notabspaces)")
	category := flag.String("category", "", "choose a pokemon from a specific category")
	name := flag.String("name", "", "choose a pokemon from a specific name")
	listCategories := flag.Bool("category-list", false, "list all available categories")

	flag.Parse()
	var args Args

	if *fastest {
		args = Args{
			Width:       *width,
			NoWrap:      true,
			TabSpaces:   "    ",
			NoTabSpaces: true,
		}
	} else {
		args = Args{
			Width:          *width,
			NoWrap:         *noWrap,
			TabSpaces:      strings.Repeat(" ", *tabWidth),
			NoTabSpaces:    *noTabSpaces,
			ListCategories: *listCategories,
			Category:       *category,
			NameToken:      *name,
		}
	}
	return args
}

func runCategoryList(categories pokedex.PokemonTrie, t *timer.Timer) {
	for k, v := range categories.Keys {
		fmt.Printf("%d (%d)\n", k, len(v))
	}
	t.Mark("ListCategories")
}

func runPrintByName(categories pokedex.PokemonTrie, args Args, t *timer.Timer) {
	matches := categories.MatchNameToken(args.NameToken)
	t.Mark("matchNameToken")
	if len(matches) == 0 {
		log.Fatal(fmt.Sprintf("Not a valid name: '%s'", args.NameToken))
	}
	printSpeechBubble(bufio.NewScanner(os.Stdin), args)
	t.Mark("printSpeechBubble")
	printPokemon(chooseRandomPokemon(matches))
	t.Mark("chooseRandomPokemon")
}

func runPrintByCategory(categories pokedex.PokemonTrie, args Args, t *timer.Timer) {
	category := []*pokedex.PokemonEntry{}
	if args.Category == "" {
		category = chooseRandomCategory(categories.Keys, categories)
		t.Mark("RandomCategory")
	} else {
		matches, err := categories.GetCategoryPaths(args.Category)
		check(err)
		category = chooseRandomCategory(matches, categories)
		t.Mark("LookupCategory")
	}

	printSpeechBubble(bufio.NewScanner(os.Stdin), args)
	t.Mark("printSpeechBubble")
	printPokemon(chooseRandomPokemon(category))
	t.Mark("chooseRandomPokemon")
}

func main() {
	args := parseFlags()
	t := timer.NewTimer()

	categories := pokedex.ReadStructFromBytes(GOBCategory)
	t.Mark("ReadCategoriesFromBytes")

	if args.ListCategories {
		runCategoryList(categories, t)
	} else if args.NameToken != "" {
		runPrintByName(categories, args, t)
	} else {
		runPrintByCategory(categories, args, t)
	}
	t.Stop()
	t.PrintJson()
}
