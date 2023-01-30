package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/tmck-code/pokesay/src/pokedex"
	"github.com/tmck-code/pokesay/src/pokesay"
)

var (
	//go:embed build/assets/pokedex.gob
	GOBCategory []byte
	//go:embed build/assets/total.txt
	GOBTotal []byte
	//go:embed build/assets/cows/*cow
	GOBCowData embed.FS
	//go:embed build/assets/metadata/*metadata
	GOBCowNames embed.FS

	MetadataRoot string = "build/assets/metadata"
	CowDataRoot  string = "build/assets/cows"
)

type Args struct {
	Width          int
	NoWrap         bool
	TabSpaces      string
	NoTabSpaces    bool
	ListCategories bool
	ListNames      bool
	Category       string
	NameToken      string
	JapaneseName   bool
}

func parseFlags() Args {
	width := flag.Int("width", 80, "the max speech bubble width")
	noWrap := flag.Bool("nowrap", false, "disable text wrapping (fastest)")
	tabWidth := flag.Int("tabwidth", 4, "replace any tab characters with N spaces")
	noTabSpaces := flag.Bool("notabspaces", false, "do not replace tab characters (fastest)")
	fastest := flag.Bool("fastest", false, "run with the fastest possible configuration (-nowrap -notabspaces)")
	category := flag.String("category", "", "choose a pokemon from a specific category")
	name := flag.String("name", "", "choose a pokemon from a specific name")
	listCategories := flag.Bool("list-categories", false, "list all available categories")
	listNames := flag.Bool("list-names", false, "list all available names")
	japaneseName := flag.Bool("japanese-name", false, "print the japanese name")

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
			ListNames:      *listNames,
			Category:       *category,
			NameToken:      *name,
			JapaneseName:   *japaneseName,
		}
	}
	return args
}

func EntryFpath(idx int) string {
	return pokedex.EntryFpath(MetadataRoot, idx)
}

func MetadataFpath(idx int) string {
	return pokedex.MetadataFpath(MetadataRoot, idx)
}

func runListCategories(categories pokedex.Trie) {
	keys := pokesay.ListCategories(categories)
	fmt.Printf("%s\n%d %s\n", strings.Join(keys, " "), len(keys), "total categories")
}

func runListNames() {
	total := pokedex.ReadIntFromBytes(GOBTotal)
	names := make([]string, total)

	for i := 0; i < total; i++ {
		metadata := pokedex.ReadMetadataFromEmbedded(GOBCowNames, MetadataFpath(i))
		names[i] = metadata.Name
	}
	fmt.Println(strings.Join(names, " "))
	fmt.Printf("\n%d %s\n", len(names), "total names")
}

func GenerateNames(metadata pokedex.PokemonMetadata, args Args) []string {
	if args.JapaneseName {
		return []string{
			metadata.Name,
			fmt.Sprintf("%s (%s)", metadata.JapaneseName, metadata.JapanesePhonetic),
		}
	} else {
		return []string{metadata.Name}
	}
}

func runPrintByName(args Args, categories pokedex.Trie) {
	match := pokesay.ChooseByName(args.NameToken, categories)
	metadata := pokedex.ReadMetadataFromEmbedded(GOBCowNames, MetadataFpath(match.Entry.Index))

	pokesay.PrintSpeechBubble(bufio.NewScanner(os.Stdin), args.Width, args.NoTabSpaces, args.TabSpaces, args.NoWrap)
	pokesay.PrintPokemon(match.Entry.Index, GenerateNames(metadata, args), match.Keys, GOBCowData)
}

func runPrintByCategory(args Args, categories pokedex.Trie) {
	choice, keys := pokesay.ChooseByCategory(args.Category, categories)
	metadata := pokedex.ReadMetadataFromEmbedded(GOBCowNames, MetadataFpath(choice.Index))

	pokesay.PrintSpeechBubble(bufio.NewScanner(os.Stdin), args.Width, args.NoTabSpaces, args.TabSpaces, args.NoWrap)
	pokesay.PrintPokemon(choice.Index, GenerateNames(metadata, args), keys, GOBCowData)
}

func runPrintRandom(args Args) {
	choice := pokesay.RandomInt(pokedex.ReadIntFromBytes(GOBTotal))

	metadata := pokedex.ReadMetadataFromEmbedded(GOBCowNames, MetadataFpath(choice))

	pokesay.PrintSpeechBubble(bufio.NewScanner(os.Stdin), args.Width, args.NoTabSpaces, args.TabSpaces, args.NoWrap)
	pokesay.PrintPokemon(choice, GenerateNames(metadata, args), strings.Split(metadata.Categories, "/"), GOBCowData)
}

func main() {
	args := parseFlags()

	if args.ListCategories {
		runListCategories(pokedex.NewTrieFromBytes(GOBCategory))
	} else if args.ListNames {
		runListNames()
	} else if args.NameToken != "" {
		runPrintByName(args, pokedex.NewTrieFromBytes(GOBCategory))
	} else if args.Category != "" {
		runPrintByCategory(args, pokedex.NewTrieFromBytes(GOBCategory))
	} else {
		runPrintRandom(args)
	}
}
