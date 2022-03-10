package main

import (
    "bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"encoding/gob"
)

type PokemonEntry struct {
	Name       string
	Data       []byte
	Categories []string
}

func findFiles(dirpath string, ext string, skip []string) []*PokemonEntry {
	pokedex := []*PokemonEntry{}
	err := filepath.Walk(dirpath, func(fpath string, f os.FileInfo, err error) error {
		for _, s := range skip {
			if strings.Contains(fpath, s) {
				return err
			}
		}
		if !f.IsDir() && filepath.Ext(f.Name()) == ext {
			data, _ := os.ReadFile(fpath)

			pokedex = append(
				pokedex,
				&PokemonEntry{Name: fpath, Data: data, },
			)
		}
		return err
	})
	if err != nil {
		fmt.Println("Fatal!")
		log.Fatal(err)
	}
	return pokedex
}

func createCategories(fpath string) []string {
	return strings.Split(fpath, "/")
}

func createPokemonEntry(pokemon *PokemonEntry, fpath string) *PokemonEntry {
	return &PokemonEntry{
		Name:       pokemon.Name,
		Data:       pokemon.Data,
		Categories: createCategories(pokemon.Name),
	}
}

func writeToFile(pokemon []*PokemonEntry, fpath string) {
    ostream, err := os.Create(fpath)
    if err != nil {
        log.Fatal(err)
    }
    writer := bufio.NewWriter(ostream)
    enc := gob.NewEncoder(writer)
    enc.Encode(pokemon)
    writer.Flush()
    ostream.Close()
	
	fmt.Println("->", fpath)
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

	pokedex := findFiles(args.FromDir, ".png", args.SkipDirs)

	writeToFile(pokedex, args.ToFpath)
}

