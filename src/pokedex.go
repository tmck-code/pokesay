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

func findFiles(dirpath string, ext string, skip []string) (pokedex.PokemonEntryMap, [][]byte) {
	categories := &pokedex.PokemonEntryMap{Categories: make(map[string][]*pokedex.PokemonEntry)}
	pokemon := make([][]byte, 0)
	idx := 0
	err := filepath.Walk(dirpath, func(fpath string, f os.FileInfo, err error) error {
		idx += 1
		for _, s := range skip {
			if strings.Contains(fpath, s) {
				return err
			}
		}
		if !f.IsDir() && filepath.Ext(f.Name()) == ext {
			data, err := os.ReadFile(fpath)
			check(err)

			pokemonCategories := createCategories(fpath)
			p := pokedex.NewPokemonEntry(
				idx,
				createName(fpath),
				tokenizeName(fpath),
			)

			for _, c := range pokemonCategories {
				if val, ok := categories.Categories[c]; ok {
					val = append(val, p)
				} else {
					categories.Categories[c] = []*pokedex.PokemonEntry{p}
				}
				categories.Categories[c] = append(categories.Categories[c], p)
			}
			pokemon = append(pokemon, data)
		}
		return err
	})
	check(err)

	return *categories, pokemon
}

func createName(fpath string) string {
	parts := strings.Split(fpath, "/")
	return strings.Split(parts[len(parts)-1], ".")[0]
}

func tokenizeName(fpath string) []string {
	return strings.Split(createName(fpath), "-")
}

func createCategories(fpath string) []string {
	parts := strings.Split(fpath, "/")
	return parts[3 : len(parts)-1]
}

type CowBuildArgs struct {
	FromDir  string
	ToFpath  string
	SkipDirs []string
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
	fmt.Println("starting at", args.FromDir)
	t := timer.NewTimer()

	categories, pokemon := findFiles(args.FromDir, ".cow", args.SkipDirs)
	t.Mark("CreateEntriesFromFiles")

	pokedex.WriteToFile(categories, args.ToFpath)
	t.Mark("WriteCategoriesToFile")

	pokedex.WriteByteToFile(pokemon, "build/data.gob")
	t.Mark("WritePokemonToFile")

	if args.DebugTimer {
		t.Stop()
		t.PrintJson()
	}
}
