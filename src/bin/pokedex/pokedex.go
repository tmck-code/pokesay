package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"

	"github.com/tmck-code/pokesay/src/bin"
	"github.com/tmck-code/pokesay/src/pokedex"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

type PokedexArgs struct {
	FromDir          string
	ToDir            string
	Debug            bool
	ToCategoryFname  string
	ToDataSubDir     string
	ToMetadataSubDir string
	ToTotalFname     string
	NewLineMode		 bool
}

func parseArgs() PokedexArgs {
	fromDir := flag.String("from", "/tmp/cows", "from dir")
	toDir := flag.String("to", "build/assets/", "to dir")

	toDataSubDir := flag.String("toDataSubDir", "cows/", "dir to write all binary (image) data to")
	toMetadataSubDir := flag.String("toMetadataSubDir", "metadata/", "dir to write all binary (metadata) data to")
	toCategoryFname := flag.String("toCategoryFpath", "pokedex.gob", "to fpath")
	toTotalFname := flag.String("toTotalFname", "total.txt", "file to write the number of available entries to")
	newlineMode := flag.Bool("newlineMode", false, "print progress via newlines, for use in CI/CD")
	debug := flag.Bool("debug", false, "show debug logs")

	flag.Parse()

	args := PokedexArgs{
		FromDir:          pokedex.NormaliseRelativeDir(*fromDir),
		ToDir:            pokedex.NormaliseRelativeDir(*toDir),
		ToCategoryFname:  *toCategoryFname,
		ToDataSubDir:     pokedex.NormaliseRelativeDir(*toDataSubDir),
		ToMetadataSubDir: pokedex.NormaliseRelativeDir(*toMetadataSubDir),
		ToTotalFname:     *toTotalFname,
		Debug:            *debug,
		NewLineMode:	  *newlineMode,
	}
	if args.Debug {
		fmt.Printf("%+v\n", args)
	}
	return args
}

func main() {
	args := parseArgs()

	totalFpath := path.Join(args.ToDir, args.ToTotalFname)
	categoryFpath := path.Join(args.ToDir, args.ToCategoryFname)

	fpaths := pokedex.FindFiles(args.FromDir, ".cow", make([]string, 0))

	err := os.MkdirAll(path.Join(args.ToDir, args.ToDataSubDir), 0755)
	check(err)
	err = os.MkdirAll(path.Join(args.ToDir, args.ToMetadataSubDir), 0755)
	check(err)

	// categories is a PokemonTrie struct that will be written to a file using encoding/gob
	// metadata is a list of pokemon data and an index to use when writing them to a file
	// - this index matches a corresponding one in the categories struct
	// - these files are embedded into the build binary using go:embed and then loaded at runtime
	categories, metadata := pokedex.CreateMetadata(args.FromDir, fpaths, args.Debug)

	pokedex.WriteStructToFile(categories, categoryFpath)

	fmt.Println("\nConverting cowfiles -> category & metadata GOB")
	pbar := bin.NewProgressBar(len(fpaths), args.NewLineMode)
	for _, m := range metadata {
		pokedex.WriteBytesToFile(m.Data, pokedex.EntryFpath(path.Join(args.ToDir, args.ToDataSubDir), m.Index), true)
		pokedex.WriteStructToFile(m.Metadata, pokedex.MetadataFpath(path.Join(args.ToDir, args.ToMetadataSubDir), m.Index))
		pbar.Add(1)
	}
	pokedex.WriteBytesToFile([]byte(strconv.Itoa(len(metadata))), totalFpath, false)

	fmt.Println("Finished converting", len(fpaths), "pokesprite -> cowfiles")
	fmt.Println("Wrote categories to", path.Join(args.ToDir, args.ToCategoryFname))
}
