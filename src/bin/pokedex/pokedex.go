package main

import (
	"flag"
	"fmt"
	"log"
	"path"
	"strconv"

	"github.com/tmck-code/pokesay-go/src/bin"
	"github.com/tmck-code/pokesay-go/src/pokedex"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

var (
	DEBUG bool = false
)

type PokedexArgs struct {
	FromDir          string
	ToDir            string
	Debug            bool
	ToCategoryFname  string
	ToDataSubDir     string
	ToMetadataSubDir string
	ToTotalFname     string
}

func parseArgs() PokedexArgs {
	fromDir := flag.String("from", "/tmp/cows", "from dir")
	toDir := flag.String("to", "build/assets/", "to dir")

	toDataSubDir := flag.String("toDataSubDir", "cows/", "dir to write all binary (image) data to")
	toMetadataSubDir := flag.String("toMetadataSubDir", "metadata/", "dir to write all binary (metadata) data to")
	toCategoryFname := flag.String("toCategoryFpath", "pokedex.gob", "to fpath")
	toTotalFname := flag.String("toTotalFname", "total.txt", "file to write the number of available entries to")
	debug := flag.Bool("debug", DEBUG, "show debug logs")

	flag.Parse()

	DEBUG = *debug
	args := PokedexArgs{
		FromDir:          *fromDir,
		ToDir:            *toDir,
		ToCategoryFname:  *toCategoryFname,
		ToDataSubDir:     *toDataSubDir,
		ToMetadataSubDir: *toMetadataSubDir,
		ToTotalFname:     *toTotalFname,
	}

	return args
}

func main() {
	args := parseArgs()

	totalFpath := path.Join(args.ToDir, args.ToTotalFname)
	categoryFpath := path.Join(args.ToDir, args.ToCategoryFname)

	fpaths := pokedex.FindFiles(args.ToDir, ".cow", make([]string, 0))

	// categories is a PokemonTrie struct that will be written to a file using encoding/gob
	// metadata is a list of pokemon data and an index to use when writing them to a file
	// - this index matches a corresponding one in the categories struct
	// - these files are embedded into the build binary using go:embed and then loaded at runtime
	categories, metadata := pokedex.CreateMetadata(fpaths)

	pokedex.WriteStructToFile(categories, categoryFpath)

	fmt.Println("\nConverting cowfiles -> category & metadata GOB")
	pbar := bin.NewProgressBar(len(fpaths))
	for _, m := range metadata {
		pokedex.WriteBytesToFile(m.Data, pokedex.EntryFpath(args.ToDataSubDir, m.Index), true)
		pokedex.WriteStructToFile(m.Metadata, pokedex.MetadataFpath(args.ToMetadataSubDir, m.Index))
		pbar.Add(1)
	}
	pokedex.WriteBytesToFile([]byte(strconv.Itoa(len(metadata))), totalFpath, false)

	fmt.Println("Finished converting", len(fpaths), "pokesprite -> cowfiles")
	fmt.Println("Wrote categories to", path.Join(args.ToDir, args.ToCategoryFname))
}
