package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/go-wordwrap"
	"github.com/tmck-code/pokesay-go/src/pokedex"
)

var (
	//go:embed build/pokedex.gob
	GOBCategory []byte
	//go:embed build/total.txt
	GOBTotal []byte
	//go:embed build/*cow
	GOBCowData embed.FS
	//go:embed build/*metadata
	GOBCowNames embed.FS

	Rand rand.Source = rand.NewSource(time.Now().UnixNano())
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

func printPokemon(index int, name string, categoryKeys []string) {
	d, _ := GOBCowData.ReadFile(pokedex.EntryFpath(index))
	fmt.Printf("%schoice: %s / categories: %s\n", pokedex.Decompress(d), name, categoryKeys)
}

func chooseRandomCategory(keys [][]string, categories pokedex.PokemonTrie) ([]string, []*pokedex.PokemonEntry) {
	choice := keys[randomInt(len(keys)-1)]
	category, err := categories.GetCategory(choice)
	check(err)
	return choice, category
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

func runCategoryList(categories pokedex.PokemonTrie) {
	ukm := map[string]bool{}
	for _, v := range categories.Keys {
		for _, k := range v {
			ukm[k] = true
		}
	}
	for k, _ := range ukm {
		fmt.Println(k)
	}
}

func runPrintByName(args Args, categories pokedex.PokemonTrie) {
	matches, err := categories.MatchNameToken(args.NameToken)
	check(err)
	match := matches[randomInt(len(matches))]

	printSpeechBubble(bufio.NewScanner(os.Stdin), args)
	printPokemon(match.Entry.Index, match.Entry.Name, match.Categories)
}

func runPrintByCategory(args Args, categories pokedex.PokemonTrie) {
	category := []*pokedex.PokemonEntry{}
	keys := []string{}

	matches, err := categories.GetCategoryPaths(args.Category)
	check(err)
	keys, category = chooseRandomCategory(matches, categories)
	choice := chooseRandomPokemon(category)

	printSpeechBubble(bufio.NewScanner(os.Stdin), args)
	printPokemon(choice.Index, choice.Name, keys)
}

func runPrintRandom(args Args) {
	total, _ := strconv.Atoi(string(GOBTotal))
	choice := randomInt(total)
	m, err := GOBCowNames.ReadFile(pokedex.MetadataFpath(choice))
	check(err)
	metadata := pokedex.ReadMetadataFromBytes(m)

	printSpeechBubble(bufio.NewScanner(os.Stdin), args)
	printPokemon(choice, metadata.Name, strings.Split(metadata.Categories, "/"))
}

func main() {
	args := parseFlags()

	if args.ListCategories {
		runCategoryList(pokedex.ReadTrieFromBytes(GOBCategory))
	} else if args.NameToken != "" {
		runPrintByName(args, pokedex.ReadTrieFromBytes(GOBCategory))
	} else if args.Category != "" {
		runPrintByCategory(args, pokedex.ReadTrieFromBytes(GOBCategory))
	} else {
		runPrintRandom(args)
	}
}
