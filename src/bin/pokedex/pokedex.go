package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/tmck-code/pokesay/src/bin"
	"github.com/tmck-code/pokesay/src/pokedex"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

// Strips the leading "./" from a path e.g. "./cows/ -> cows/"
func normaliseRelativeDir(dirPath string) string {
	return strings.TrimPrefix(dirPath, "./")
}

type PokedexArgs struct {
	FromDir           string
	FromMetadataFname string
	ToDir             string
	Debug             bool
	ToCategoryFname   string
	ToDataSubDir      string
	ToMetadataSubDir  string
	ToTotalFname      string
}

func parseArgs() PokedexArgs {
	fromDir := flag.String("from", "/tmp/cows", "from dir")
	fromMetadataFname := flag.String("fromMetadata", "/tmp/cows/pokemon.json", "metadata file")
	toDir := flag.String("to", "build/assets/", "to dir")

	toDataSubDir := flag.String("toDataSubDir", "cows/", "dir to write all binary (image) data to")
	toMetadataSubDir := flag.String("toMetadataSubDir", "metadata/", "dir to write all binary (metadata) data to")
	toCategoryFname := flag.String("toCategoryFpath", "pokedex.gob", "to fpath")
	toTotalFname := flag.String("toTotalFname", "total.txt", "file to write the number of available entries to")
	debug := flag.Bool("debug", false, "show debug logs")

	flag.Parse()

	args := PokedexArgs{
		FromDir:           normaliseRelativeDir(*fromDir),
		FromMetadataFname: *fromMetadataFname,
		ToDir:             normaliseRelativeDir(*toDir),
		ToCategoryFname:   *toCategoryFname,
		ToDataSubDir:      normaliseRelativeDir(*toDataSubDir),
		ToMetadataSubDir:  normaliseRelativeDir(*toMetadataSubDir),
		ToTotalFname:      *toTotalFname,
		Debug:             *debug,
	}
	if args.Debug {
		fmt.Printf("%+v\n", args)
	}
	return args
}

// This function reads in the files given by the PokedexArgs, and generates the data that pokesay will use when running
// - The "category" struct
//   - contains category information, and the index of the corresponding metadata file
//
// - The "metadata" files
//   - named like 1.metadata, contains pokemon info like name, categories, japanese name
//
// - The "data" files
//   - contain the pokemon as gzipped text
//
// - The "total" file
//   - contains the total number of pokemon files, used for random selection
func main() {
	args := parseArgs()

	totalFpath := path.Join(args.ToDir, args.ToTotalFname)
	categoryFpath := path.Join(args.ToDir, args.ToCategoryFname)
	entryDirPath := path.Join(args.ToDir, args.ToDataSubDir)
	metadataDirPath := path.Join(args.ToDir, args.ToMetadataSubDir)

	cowfileFpaths := pokedex.FindFiles(args.FromDir, ".cow", make([]string, 0))

	err := os.MkdirAll(entryDirPath, 0755)
	check(err)
	err = os.MkdirAll(metadataDirPath, 0755)
	check(err)

	pokemonNames := pokedex.ReadNames(args.FromMetadataFname)
	fmt.Printf("%+v\n", pokemonNames)

	// 1. Create the category struct using the cowfile paths, pokemon names and indexes
	categories := pokedex.CreateCategoryStruct(args.FromDir, cowfileFpaths, args.Debug)

	// 2. Create the metadata files, containing name/category/japanese name info for each pokemon
	metadata := pokedex.CreateMetadata(args.FromDir, cowfileFpaths, pokemonNames, args.Debug)

	// categories is a Trie struct that will be written to a file using encoding/gob
	// metadata is a list of pokemon data and an index to use when writing them to a file
	// - this index matches a corresponding one in the categories struct
	// - these files are embedded into the build binary using go:embed and then loaded at runtime
	// categories, metadata := pokedex.CreateMetadata(args.FromDir, cowfileFpaths, pokemonNames, args.Debug)

	pokedex.WriteStructToFile(categories, categoryFpath)

	fmt.Println("\nConverting cowfiles -> category & metadata GOB")
	pbar := bin.NewProgressBar(len(cowfileFpaths))
	for _, m := range metadata {
		pokedex.WriteBytesToFile(m.Data, pokedex.EntryFpath(entryDirPath, m.Index), true)
		pokedex.WriteStructToFile(m.Metadata, pokedex.MetadataFpath(metadataDirPath, m.Index))
		pbar.Add(1)
	}
	pokedex.WriteBytesToFile([]byte(strconv.Itoa(len(metadata))), totalFpath, false)

	fmt.Println("Finished converting", len(cowfileFpaths), "pokesprite -> cowfiles")
	fmt.Println("Wrote categories to", path.Join(args.ToDir, args.ToCategoryFname))
}
