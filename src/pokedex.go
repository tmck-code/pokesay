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
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func findFiles(dirpath string, ext string, skip []string) pokedex.PokemonEntryMap {
	categories := &pokedex.PokemonEntryMap{Categories: make(map[string][]pokedex.PokemonEntry)}
	err := filepath.Walk(dirpath, func(fpath string, f os.FileInfo, err error) error {
		for _, s := range skip {
			if strings.Contains(fpath, s) {
				return err
			}
		}
		if !f.IsDir() && filepath.Ext(f.Name()) == ext {
			data, err := os.ReadFile(fpath)
			check(err)
			pokemonCategories := createCategories(fpath)

			for _, c := range pokemonCategories {
				p := pokedex.NewPokemonEntry(
					data,
					createName(fpath),
					tokenizeName(fpath),
					createCategories(fpath),
				)
				if val, ok := categories.Categories[c]; ok {
					val = append(val, *p)
				} else {
					categories.Categories[c] = []pokedex.PokemonEntry{*p}
				}
				categories.Categories[c] = append(categories.Categories[c], *p)
			}
		}
		return err
	})
	check(err)

	categories.NCategories = len(categories.Categories)

	return *categories
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
	return parts[2 : len(parts)-1]
}

type CowBuildArgs struct {
	FromDir  string
	ToFpath  string
	SkipDirs []string
}

func parseArgs() CowBuildArgs {
	fromDir := flag.String("from", ".", "from dir")
	toFpath := flag.String("to", ".", "to fpath")
	skipDirs := flag.String("skip", "'[\"resources\"]'", "JSON array of dir patterns to skip converting")

	flag.Parse()

	args := CowBuildArgs{FromDir: *fromDir, ToFpath: *toFpath}
	json.Unmarshal([]byte(*skipDirs), &args.SkipDirs)

	return args
}

func main() {
	args := parseArgs()
	fmt.Println("starting at", args.FromDir)

	categories := findFiles(args.FromDir, ".cow", args.SkipDirs)

	pokedex.WriteToFile(categories, args.ToFpath)
}
