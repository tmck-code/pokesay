package main

import (
	"embed"
	"flag"
	"fmt"

	"github.com/tmck-code/pokesay/src/pokedex"
)

var (
	//go:embed assets/metadata/*metadata
	GOBCowNames embed.FS
	//go:embed assets/cows/*cow
	GOBCowFiles  embed.FS
	MetadataRoot string = "assets/metadata"
	CowfileRoot  string = "assets/cows"
)

type Args struct {
	Index int
}

func parseFlags() Args {
	index := flag.Int("index", 80, "the metadata file index")
	flag.Parse()

	return Args{
		Index: *index,
	}
}
func MetadataFpath(idx int) string {
	return pokedex.MetadataFpath(MetadataRoot, idx)
}

func EntryFpath(idx int) string {
	return pokedex.EntryFpath(CowfileRoot, idx)
}

func main() {
	args := parseFlags()
	metadata := pokedex.ReadMetadataFromEmbedded(GOBCowNames, MetadataFpath(args.Index))
	fmt.Println(pokedex.StructToJSON(metadata, 2))

	for _, entry := range metadata.Entries {
		fmt.Println(
			entry.Categories,
			"\n",
			string(pokedex.ReadPokemonCow(GOBCowFiles, EntryFpath(entry.EntryIndex))),
		)
	}
}
