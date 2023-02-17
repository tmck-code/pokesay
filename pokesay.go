package main

import (
	"embed"
	"flag"
	"fmt"
	"path"
	"strconv"
	"strings"

	"github.com/tmck-code/pokesay/src/pokedex"
	"github.com/tmck-code/pokesay/src/pokesay"
	"github.com/tmck-code/pokesay/src/timer"
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
	//go:embed all:build/assets/categories
	GOBCategories embed.FS
	//go:embed build/assets/categories.txt
	GOBCategoryKeys []byte

	MetadataRoot string = "build/assets/metadata"
	CowDataRoot  string = "build/assets/cows"
)

func parseFlags() pokesay.Args {
	width := flag.Int("width", 80, "the max speech bubble width")
	noWrap := flag.Bool("no-wrap", false, "disable text wrapping (fastest)")
	tabWidth := flag.Int("tab-width", 4, "replace any tab characters with N spaces")
	noTabSpaces := flag.Bool("no-tab-spaces", false, "do not replace tab characters (fastest)")
	noCategoryInfo := flag.Bool("no-category-info", false, "do not print pokemon categories")
	fastest := flag.Bool("fastest", false, "run with the fastest possible configuration (-nowrap -notabspaces)")
	category := flag.String("category", "", "choose a pokemon from a specific category")
	name := flag.String("name", "", "choose a pokemon from a specific name")
	listCategories := flag.Bool("list-categories", false, "list all available categories")
	listNames := flag.Bool("list-names", false, "list all available names")
	japaneseName := flag.Bool("japanese-name", false, "print the japanese name")

	flag.Parse()
	var args pokesay.Args

	if *fastest {
		args = pokesay.Args{
			Width:       *width,
			NoWrap:      true,
			TabSpaces:   "    ",
			NoTabSpaces: true,
		}
	} else {
		args = pokesay.Args{
			Width:          *width,
			NoWrap:         *noWrap,
			TabSpaces:      strings.Repeat(" ", *tabWidth),
			NoTabSpaces:    *noTabSpaces,
			NoCategoryInfo: *noCategoryInfo,
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
	return pokedex.EntryFpath(CowDataRoot, idx)
}

func MetadataFpath(idx int) string {
	return pokedex.MetadataFpath(MetadataRoot, idx)
}

func runListCategories() {
	keys := pokedex.ReadStructFromBytes[[]string](GOBCategoryKeys)
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

func GenerateNames(metadata pokedex.PokemonMetadata, args pokesay.Args) []string {
	if args.JapaneseName {
		return []string{
			metadata.Name,
			fmt.Sprintf("%s (%s)", metadata.JapaneseName, metadata.JapanesePhonetic),
		}
	} else {
		return []string{metadata.Name}
	}
}

func runPrintByName(args pokesay.Args, categories pokedex.Trie) {
	t := timer.NewTimer("runPrintByName", true)
	match := pokesay.ChooseByName(args.NameToken, categories)
	t.Mark("match")
	metadata := pokedex.ReadMetadataFromEmbedded(GOBCowNames, MetadataFpath(match.Entry.Index))
	t.Mark("read metadata")
	choice := pokesay.RandomInt(len(metadata.Entries))
	t.Mark("choice")
	final := metadata.Entries[choice]
	// fmt.Println(pokedex.StructToJSON(metadata), "\n", choice, "\n", final)

	pokesay.Print(args, final.EntryIndex, GenerateNames(metadata, args), final.Categories, GOBCowData)
	t.Mark("print")

	t.Stop()
	t.PrintJson()
}

func runPrintByCategory(args pokesay.Args) {
	t := timer.NewTimer("runPrintByCategory", true)

	dirPath := fmt.Sprintf("build/assets/categories/%s", args.Category)
	dir, _ := GOBCategories.ReadDir(dirPath)
	// fmt.Printf("----- %d %#v\n", len(dir), dir)

	choice := dir[pokesay.RandomInt(len(dir))]
	// fmt.Printf("random category choice %#v\n", choice)

	// fmt.Println("category data", choice)

	categoryMetadata, err := GOBCategories.ReadFile(fmt.Sprintf("build/assets/categories/%s/%s", args.Category, choice.Name()))
	pokesay.Check(err)

	parts := strings.Split(string(categoryMetadata), "/")
	// fmt.Println("parts:", parts)
	metadataIndex, err := strconv.Atoi(string(parts[0]))
	pokesay.Check(err)

	entryIndex, err := strconv.Atoi(string(parts[1]))
	pokesay.Check(err)

	// fmt.Printf("category metadata: %#v\n", categoryMetadata)

	metadata := pokedex.ReadMetadataFromEmbedded(GOBCowNames, MetadataFpath(metadataIndex))
	// fmt.Printf("name metadata: %#v\n", metadata)

	t.Mark("read file")
	t.Mark("read metadata")
	// fmt.Println(entryIndex, metadata.Entries)
	final := metadata.Entries[entryIndex]
	t.Mark("choose entry")

	pokesay.Print(args, final.EntryIndex, GenerateNames(metadata, args), final.Categories, GOBCowData)
	t.Mark("print")

	t.Stop()
	t.PrintJson()
}

func runPrintRandom(args pokesay.Args) {
	t := timer.NewTimer("runPrintRandom", true)
	_, choice := pokesay.ChooseByRandomIndex(GOBTotal)
	t.Mark("choose index")
	metadata := pokedex.ReadMetadataFromEmbedded(GOBCowNames, MetadataFpath(choice))
	t.Mark("read metadata")

	final := metadata.Entries[pokesay.RandomInt(len(metadata.Entries))]
	t.Mark("choose entry")

	pokesay.Print(args, final.EntryIndex, GenerateNames(metadata, args), final.Categories, GOBCowData)
	t.Mark("print")

	t.Stop()
	t.PrintJson()
}

func main() {
	args := parseFlags()
	t := timer.NewTimer("main", true)

	if args.ListCategories {
		runListCategories()
		t.Mark("op")
	} else if args.ListNames {
		runListNames()
	} else if args.NameToken != "" {
		c := pokedex.NewTrieFromBytes(GOBCategory)
		t.Mark("trie")
		runPrintByName(args, c)
		t.Mark("op")
	} else if args.Category != "" {
		runPrintByCategory(args)
		t.Mark("op")
	} else {
		runPrintRandom(args)
		t.Mark("op")
	}

	t.Stop()
	t.PrintJson()
}
