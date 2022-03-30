package main

import (
	"bufio"
	"embed"
	_ "encoding/binary"
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
	//go:embed build/cows.gob
	GOBCategory []byte
	//go:embed build/*cow
	GOBCowData embed.FS
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
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

func randomInt(n int) int {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Intn(n)
}

func printPokemon(choice pokedex.PokemonEntry) {
	d, _ := GOBCowData.ReadFile(pokedex.EntryFpath(choice.Index))
	fmt.Printf("%s\nchoice: %s\n", pokedex.Decompress(d), choice.Name, choice.Index)
}

func chooseRandomCategory(entries pokedex.PokemonTrie) []*pokedex.PokemonEntry {
	categoryChoice := entries.Keys[randomInt(len(entries.Keys)-1)]
	category, ok := entries.GetCategory(categoryChoice)
	if !ok {
		log.Fatal("Couldn't choose random category")
	}
	return category
}

func chooseRandomPokemon(pokemon []*pokedex.PokemonEntry) pokedex.PokemonEntry {
	return *pokemon[randomInt(len(pokemon))]
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
		fmt.Printf("%s (%d)\n", k, len(v))
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
	category := []string{}
	if args.Category == "" {
		category = categories.Keys[randomInt(len(categories.Keys)-1)]
		t.Mark("RandomCategory")
	} else {
		category = categories.GetCategoryPath(args.Category)
		t.Mark("LookupCategory")
		if len(category) == 0 {
			log.Fatal(fmt.Sprintf("Category not found: %s", category))
		}
	}

	printSpeechBubble(bufio.NewScanner(os.Stdin), args)
	t.Mark("printSpeechBubble")
	matches, ok := categories.GetCategory(category)
	t.Mark("GetCategory")
	if !ok {
		log.Fatal(fmt.Sprintf("No pokemon found for category: %s", category))
	}
	printPokemon(chooseRandomPokemon(matches))
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
