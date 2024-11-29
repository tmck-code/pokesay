package main

import (
	"embed"
	"fmt"
	"strings"

	"github.com/pborman/getopt/v2"
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

	CategoryRoot string = "build/assets/categories" // the root directory of the pokemon categories
	MetadataRoot string = "build/assets/metadata"   // the root directory of the pokemon metadata
	CowDataRoot  string = "build/assets/cows"       // the root directory of the pokemon cow data
)

// parseFlags parses the command line flags and returns a pokesay.Args struct
func parseFlags() pokesay.Args {
	help := getopt.BoolLong("help", 'h', "display this help message")
	// print verbose output (currently timer output)
	verbose := getopt.BoolLong("verbose", 'v', "print verbose output", "verbose")

	// selection/filtering
	name := getopt.StringLong("name", 'n', "", "choose a pokemon from a specific name")
	category := getopt.StringLong("category", 'c', "", "choose a pokemon from a specific category")

	// list operations
	listNames := getopt.BoolLong("list-names", 'l', "list all available names")
	listCategories := getopt.BoolLong("list-categories", 'L', "list all available categories")

	width := getopt.IntLong("width", 'w', 80, "the max speech bubble width")

	// speech bubble options
	tabWidth := getopt.IntLong("tab-width", 't', 4, "replace any tab characters with N spaces")
	noWrap := getopt.BoolLong("no-wrap", 'W', "disable text wrapping (fastest)")
	noTabSpaces := getopt.BoolLong("no-tab-spaces", 's', "do not replace tab characters (fastest)")
	fastest := getopt.BoolLong("fastest", 'f', "run with the fastest possible configuration (--nowrap & --notabspaces)")
	noBubble := getopt.BoolLong("no-bubble", 'B', "do not draw the speech bubble")

	// info box options
	japaneseName := getopt.BoolLong("japanese-name", 'j', "print the japanese name in the info box")
	noCategoryInfo := getopt.BoolLong("no-category-info", 'C', "do not print pokemon category information in the info box")
	drawInfoBorder := getopt.BoolLong("info-border", 'b', "draw a border around the info box")

	// other option
	unicodeBorders := getopt.BoolLong("unicode-borders", 'u', "use unicode characters to draw the border around the speech box (and info box if --info-border is enabled)")

	getopt.Parse()
	var args pokesay.Args

	if *fastest {
		args = pokesay.Args{
			Width:       *width,
			NoWrap:      true,
			TabSpaces:   "    ",
			NoTabSpaces: true,
			BoxChars:    pokesay.DetermineBoxChars(false),
			Help:        *help,
			Verbose:     *verbose,
		}
	} else {
		args = pokesay.Args{
			Width:          *width,
			NoWrap:         *noWrap,
			DrawBubble:     !*noBubble,
			TabSpaces:      strings.Repeat(" ", *tabWidth),
			NoTabSpaces:    *noTabSpaces,
			NoCategoryInfo: *noCategoryInfo,
			ListCategories: *listCategories,
			ListNames:      *listNames,
			Category:       *category,
			NameToken:      *name,
			JapaneseName:   *japaneseName,
			BoxChars:       pokesay.DetermineBoxChars(*unicodeBorders),
			DrawInfoBorder: *drawInfoBorder,
			Help:           *help,
			Verbose:        *verbose,
		}
	}
	return args
}

// runListCategories prints all available categories
// - This reads a list of categories from the embedded filesystem
// - prints the list of categories, and the total number of categories
func runListCategories() {
	categories := pokedex.ReadStructFromBytes[[]string](GOBCategoryKeys)
	fmt.Printf("%s\n%d %s\n", strings.Join(categories, " "), len(categories), "total categories")
}

// runListNames prints all available pokemon names
// - This reads a struct of {name -> metadata indexes} from the embedded filesystem
// - prints all the keys of the struct, and the total number of names
func runListNames() {
	names := pokesay.ListNames(
		pokedex.ReadStructFromBytes[map[string][]int](GOBAllNames),
	)
	fmt.Printf("%s\n%d %s\n", strings.Join(names, " "), len(names), "total names")
}

// GenerateNames returns a list of names to print
// - If the japanese name flag is set, it returns both the english and japanese names
// - Otherwise, it returns just the english name
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

// runPrintByName prints a pokemon matched by a name
// The name must match the lowercase name of the pokemon (TODO: improve this behaviour)
// - This reads a struct of {name -> metadata indexes} from the embedded filesystem
// - It matches the name to a metadata index, loads the corresponding metadata file, and then chooses a random entry
// - Finally, it prints the pokemon
func runPrintByName(args pokesay.Args) {
	t := timer.NewTimer("runPrintByName", true)

	names := pokedex.ReadStructFromBytes[map[string][]int](GOBAllNames)
	t.Mark("read name struct")

	metadata, final := pokesay.ChooseByName(names, args.NameToken, GOBCowNames, MetadataRoot)
	t.Mark("find/read metadata")

	pokesay.Print(args, final.EntryIndex, GenerateNames(metadata, args), final.Categories, GOBCowData)
	t.Mark("print")

	t.Stop()
	t.PrintJson()
}

// runPrintByCategory prints a pokemon matched by a category
// - This loads a GOB file containing a pokemon "category" search struct from the embedded filesystem
// - It chooses a random category file from the corresponding category directory
// - It reads the category file and chooses a random pokemon from the category
//   - # TODO: as each category dir contains files with a singular entry of {metadata index/entry index},
//   - # this means that pokemon that are in the same category multiple times will be chosen more often
//
// - It reads the metadata file of the chosen pokemon and chooses the corresponding entry from the category search
// - Finally, it prints the pokemon
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

// runPrintByNameAndCategory prints a pokemon matched by a name and category
// - This reads a struct of {name -> metadata indexes} from the embedded filesystem
// - It matches the name to a metadata index, loads the corresponding metadata file, and then randomly chooses an entry that matches the category
// - Finally, it prints the pokemon
func runPrintByNameAndCategory(args pokesay.Args) {
	t := timer.NewTimer("runPrintByNameAndCategory", true)

	names := pokedex.ReadStructFromBytes[map[string][]int](GOBAllNames)
	t.Mark("read name struct")

	metadata, final := pokesay.ChooseByNameAndCategory(names, args.NameToken, GOBCowNames, MetadataRoot, args.Category)
	t.Mark("find/read metadata")

	pokesay.Print(args, final.EntryIndex, GenerateNames(metadata, args), final.Categories, GOBCowData)
	t.Mark("print")

	t.Stop()
	t.PrintJson()
}

// runPrintRandom prints a random pokemon
// - This loads a specific GOB file from the embedded filesystem that contains the number of pokemon
// - generates a random number between 0 and the number of pokemon
// - reads the metadata file of at `<index>.metadata` as a PokemonMetadata struct
// - chooses a random entry from the metadata file
// - finally prints the pokemon
func runPrintRandom(args pokesay.Args) {
	t := timer.NewTimer("runPrintRandom", true)
	choice := pokesay.RandomInt(pokedex.ReadIntFromBytes(GOBTotal))
	t.Mark("choose index")
	metadata := pokedex.ReadMetadataFromEmbedded(GOBCowNames, pokedex.MetadataFpath(MetadataRoot, choice))
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
	// if the -h/--help flag is set, print usage and exit
	if args.Help {
		getopt.Usage()
		return
	}
	if args.Verbose {
		fmt.Println("Verbose output enabled")
		timer.DEBUG = true
	}

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
