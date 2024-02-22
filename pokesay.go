package main

import (
	"embed"
	"flag"
	"fmt"
	"strings"

	"github.com/tmck-code/pokesay/src/pokedex"
	"github.com/tmck-code/pokesay/src/pokesay"
	"github.com/tmck-code/pokesay/src/timer"
)

var (
	//go:embed build/assets/category_keys.txt
	GOBCategoryKeys []byte
	//go:embed build/assets/names.txt
	GOBAllNames []byte

	//go:embed build/assets/total.txt
	GOBTotal []byte
	//go:embed build/assets/cows/*cow
	GOBCowData embed.FS
	//go:embed build/assets/metadata/*metadata
	GOBCowNames embed.FS
	//go:embed all:build/assets/categories
	GOBCategories embed.FS

	CategoryRoot string = "build/assets/categories"
	MetadataRoot string = "build/assets/metadata"
	CowDataRoot  string = "build/assets/cows"
)

func parseFlags() pokesay.Args {
	// list operations
	listCategories := flag.Bool("list-categories", false, "list all available categories")
	listNames := flag.Bool("list-names", false, "list all available names")

	// selection/filtering
	category := flag.String("category", "", "choose a pokemon from a specific category")
	name := flag.String("name", "", "choose a pokemon from a specific name")
	width := flag.Int("width", 80, "the max speech bubble width")

	// speech bubble options
	noWrap := flag.Bool("no-wrap", false, "disable text wrapping (fastest)")
	tabWidth := flag.Int("tab-width", 4, "replace any tab characters with N spaces")
	noTabSpaces := flag.Bool("no-tab-spaces", false, "do not replace tab characters (fastest)")
	fastest := flag.Bool("fastest", false, "run with the fastest possible configuration (-nowrap -notabspaces)")

	// info box options
	japaneseName := flag.Bool("japanese-name", false, "print the japanese name")
	noCategoryInfo := flag.Bool("no-category-info", false, "do not print pokemon categories")
	drawInfoBorder := flag.Bool("info-border", false, "draw a border around the info line")

	// other option
	unicodeBorders := flag.Bool("unicode-borders", false, "use unicode characters to draw the border around the speech box (and info box if -info-border is enabled)")

	flag.Parse()
	var args pokesay.Args

	if *fastest {
		args = pokesay.Args{
			Width:         *width,
			NoWrap:        true,
			TabSpaces:     "    ",
			NoTabSpaces:   true,
			BoxCharacters: pokesay.DetermineBoxCharacters(false),
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
			BoxCharacters:  pokesay.DetermineBoxCharacters(*unicodeBorders),
			DrawInfoBorder: *drawInfoBorder,
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

func CategoryFpath(category string, fname string) string {
	return pokedex.CategoryFpath(CategoryRoot, category, fname)
}

func runListCategories() {
	categories := pokedex.ReadStructFromBytes[[]string](GOBCategoryKeys)
	fmt.Printf("%s\n%d %s\n", strings.Join(categories, " "), len(categories), "total categories")
}

func runListNames() {
	names := pokesay.ListNames(
		pokedex.ReadStructFromBytes[map[string][]int](GOBAllNames),
	)
	fmt.Printf("%s\n%d %s\n", strings.Join(names, " "), len(names), "total names")
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

func runPrintByName(args pokesay.Args) {
	t := timer.NewTimer("runPrintByName", true)

	names := pokedex.ReadStructFromBytes[map[string][]int](GOBAllNames)
	t.Mark("read name struct")

	metadata, final := pokesay.ChooseByName(names, args.NameToken, GOBCowNames, MetadataRoot, "")
	t.Mark("find and read metadata")

	pokesay.Print(args, final.EntryIndex, GenerateNames(metadata, args), final.Categories, GOBCowData)
	t.Mark("print")

	t.Stop()
	t.PrintJson()
}

func runPrintByCategory(args pokesay.Args) {
	t := timer.NewTimer("runPrintByCategory", true)

	dirPath := pokedex.CategoryDirpath(CategoryRoot, args.Category)
	dir, _ := GOBCategories.ReadDir(dirPath)
	metadata, final := pokesay.ChooseByCategory(args.Category, dir, GOBCategories, CategoryRoot, GOBCowNames, MetadataRoot)

	pokesay.Print(args, final.EntryIndex, GenerateNames(metadata, args), final.Categories, GOBCowData)
	t.Mark("print")

	t.Stop()
	t.PrintJson()
}

func runPrintByNameAndCategory(args pokesay.Args) {
	t := timer.NewTimer("runPrintByNameAndCategory", true)

	names := pokedex.ReadStructFromBytes[map[string][]int](GOBAllNames)
	t.Mark("read name struct")

	metadata, final := pokesay.ChooseByName(names, args.NameToken, GOBCowNames, MetadataRoot, args.Category)
	t.Mark("find and read metadata")

	pokesay.Print(args, final.EntryIndex, GenerateNames(metadata, args), final.Categories, GOBCowData)
	t.Mark("print")

	t.Stop()
	t.PrintJson()
}

func runPrintRandom(args pokesay.Args) {
	t := timer.NewTimer("runPrintRandom", true)
	choice := pokesay.RandomInt(pokedex.ReadIntFromBytes(GOBTotal))
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
	} else if args.ListNames {
		runListNames()
	} else if args.NameToken != "" && args.Category != "" {
		runPrintByNameAndCategory(args)
	} else if args.NameToken != "" {
		runPrintByName(args)
	} else if args.Category != "" {
		runPrintByCategory(args)
	} else {
		runPrintRandom(args)
	}
	t.Mark("op")

	t.Stop()
	t.PrintJson()
}
