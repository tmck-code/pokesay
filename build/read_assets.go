package main

import (
	"embed"
	"flag"
	"fmt"
	"path/filepath"

	"github.com/tmck-code/pokesay/src/pokedex"
	"github.com/tmck-code/pokesay/src/timer"
)

var (
	//go:embed assets/metadata/*metadata
	GOBCowNames embed.FS
	//go:embed assets/cows/*cow
	GOBCowFiles embed.FS

	MetadataRoot string = "assets/metadata"
	CowfileRoot  string = "assets/cows"
)

type Args struct {
	Index int
	Fpath string
}

func parseFlags() Args {
	index := flag.Int("index", 80, "the metadata file index")
	fpath := flag.String("fpath", "", "choose a pokemon from a specific fpath")
	flag.Parse()

	return Args{
		Index: *index,
		Fpath: *fpath,
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
	t := timer.NewTimer("read_assets", true)

	var metadata pokedex.PokemonMetadata

	if args.Fpath == "" {
		fpath := MetadataFpath(args.Index)
		metadata = pokedex.ReadMetadataFromEmbedded(GOBCowNames, fpath)
	} else {
		fpath, _ := filepath.Abs(args.Fpath)
		metadata = pokedex.ReadMetadataFromFile(fpath)
	}

	t.Mark("metadata")

	fmt.Println(pokedex.StructToJSON(metadata, 2))
	t.Mark("toJSON")

	for i, entry := range metadata.Entries {
		data := string(pokedex.ReadPokemonCow(GOBCowFiles, EntryFpath(entry.EntryIndex)))
		t.Mark(fmt.Sprintf("read-cow-%d", i))

		fmt.Printf("%s\n%s\n", entry.Categories, data)
		t.Mark(fmt.Sprintf("print-cow-%d", i))
	}
	t.Stop()
	t.PrintJson()
}
