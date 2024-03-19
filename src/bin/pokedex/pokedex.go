package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/tmck-code/pokesay/src/bin"
	"github.com/tmck-code/pokesay/src/pokedex"
)

// Strips the leading "./" from a path e.g. "./cows/ -> cows/"
func normaliseRelativeDir(dirPath string) string {
	return strings.TrimPrefix(dirPath, "./")
}

type PokedexArgs struct {
	FromDir           string
	FromMetadataFname string
	ToDir             string
	Debug             bool
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
	toTotalFname := flag.String("toTotalFname", "total.txt", "file to write the number of available entries to")
	debug := flag.Bool("debug", false, "show debug logs")

	flag.Parse()

	args := PokedexArgs{
		FromDir:           normaliseRelativeDir(*fromDir),
		FromMetadataFname: *fromMetadataFname,
		ToDir:             normaliseRelativeDir(*toDir),
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

func mkDirs(dirPaths []string) {
	for _, dirPath := range dirPaths {
		err := os.MkdirAll(dirPath, 0755)
		pokedex.Check(err)
	}
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
	entryDirPath := path.Join(args.ToDir, args.ToDataSubDir)
	metadataDirPath := path.Join(args.ToDir, args.ToMetadataSubDir)

	// ensure that the destination directories exist
	mkDirs([]string{entryDirPath, metadataDirPath})

	// Find all the cowfiles
	cowfileFpaths := pokedex.FindFiles(args.FromDir, ".cow", make([]string, 0))
	fmt.Println("- Found", len(cowfileFpaths), "cowfiles")
	// Read pokemon names
	pokemonNames := pokedex.ReadNames(args.FromMetadataFname)
	fmt.Println("- Read", len(pokemonNames), "pokemon names from", args.FromMetadataFname)

	fmt.Println("- Writing entries to file")
	pbar := bin.NewProgressBar(len(cowfileFpaths))
	for i, fpath := range cowfileFpaths {
		data, err := os.ReadFile(fpath)
		pokedex.Check(err)

		pokedex.WriteBytesToFile(data, pokedex.EntryFpath(entryDirPath, i), true)
		pbar.Add(1)
	}

	// 1. For each pokemon name, write a metadata file, containing the name information, and
	// links to all of the matching cowfile indexes
	fmt.Println("- Writing metadata to file")
	pokemonMetadata := make([]pokedex.PokemonMetadata, 0)
	uniqueNames := make(map[string][]int)
	i := 0
	pbar = bin.NewProgressBar(len(pokemonNames))
	for key, name := range pokemonNames {
		metadata := pokedex.CreateNameMetadata(i, key, name, args.FromDir, cowfileFpaths)
		pokedex.WriteStructToFile(metadata, pokedex.MetadataFpath(metadataDirPath, i))
		pokemonMetadata = append(pokemonMetadata, *metadata)
		uniqueNames[name.Slug] = append(uniqueNames[name.Slug], i)
		i++
		pbar.Add(1)
	}

	pokedex.WriteStructToFile(uniqueNames, "build/assets/names.txt")

	// 2. Create the category struct using the cowfile paths, pokemon names and indexes\
	fmt.Println("- Writing categories to file")
	categories := pokedex.CreateCategoryStruct(args.FromDir, pokemonMetadata, args.Debug)
	pokedex.WriteStructToFile(categories, "build/assets/category_keys.txt")

	fmt.Println("- Writing total metadata to file")
	pokedex.WriteIntToFile(len(pokemonMetadata), totalFpath)

	fmt.Println("✓ Complete! Indexed", len(cowfileFpaths), "total cowfiles")
	fmt.Println("wrote", i, "names to", "build/assets/names.txt")

	fmt.Println("✓ Wrote gzipped metadata to", metadataDirPath)
	fmt.Println("✓ Wrote gzipped cowfiles to", entryDirPath)
	fmt.Println("✓ Wrote 'total' metadata to", totalFpath, len(pokemonMetadata))
}
