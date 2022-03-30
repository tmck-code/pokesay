package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/tmck-code/pokesay-go/src/pokedex"
	"github.com/tmck-code/pokesay-go/src/timer"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

type Metadata struct {
	Data  []byte
	Index int
}

func findFiles(dirpath string, ext string, skip []string) (pokedex.PokemonTrie, []Metadata) {
	categories := pokedex.NewTrie()
	metadata := []Metadata{}

	idx := 0
	err := filepath.Walk(dirpath, func(fpath string, f os.FileInfo, err error) error {
		for _, s := range skip {
			if strings.Contains(fpath, s) {
				return err
			}
		}
		if !f.IsDir() && filepath.Ext(f.Name()) == ext {
			idx += 1
			data, err := os.ReadFile(fpath)
			check(err)

			categories.Insert(
				createCategories(fpath),
				pokedex.NewPokemonEntry(idx, createName(fpath)),
			)
			metadata = append(metadata, Metadata{data, idx})
		}
		return err
	})
	check(err)
	fmt.Println("Wrote", idx, "pokemon to file")

	return *categories, metadata
}

func createName(fpath string) string {
	parts := strings.Split(fpath, "/")
	return strings.Split(parts[len(parts)-1], ".")[0]
}

func createCategories(fpath string) []string {
	parts := strings.Split(fpath, "/")
	return append([]string{"pokemon"}, parts[3:len(parts)-1]...)
}

type CowBuildArgs struct {
	FromDir    string
	ToFpath    string
	SkipDirs   []string
	DebugTimer bool
}

func parseArgs() CowBuildArgs {
	fromDir := flag.String("from", ".", "from dir")
	toFpath := flag.String("to", ".", "to fpath")
	skipDirs := flag.String("skip", "'[\"resources\"]'", "JSON array of dir patterns to skip converting")
	debugTimer := flag.Bool("debugTimer", false, "show a debug timer")

	flag.Parse()

	args := CowBuildArgs{FromDir: *fromDir, ToFpath: *toFpath, DebugTimer: *debugTimer}
	json.Unmarshal([]byte(*skipDirs), &args.SkipDirs)

	return args
}

func main() {
	args := parseArgs()
	t := timer.NewTimer()
	fmt.Println("starting at", args.FromDir)

	// categories is a PokemonTrie struct that will be written to a file using encoding/gob
	// metadata is a list of pokemon data and an index to use when writing them to a file
	// - this index matches a corresponding one in the categories struct
	// - these files are embedded into the build binary using go:embed and then loaded at runtime
	categories, metadata := findFiles(args.FromDir, ".cow", args.SkipDirs)
	t.Mark("CreateEntriesFromFiles")

	pokedex.WriteStructToFile(categories, args.ToFpath)
	t.Mark("WriteCategoriesToFile")

	for _, m := range metadata {
		pokedex.WriteCompressedToFile(m.Data, pokedex.EntryFpath(m.Index))
	}
	t.Mark("WriteDataToFiles")

	if args.DebugTimer {
		t.Stop()
		t.PrintJson()
	}
}
