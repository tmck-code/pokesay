package main

import (
	"bufio"
	"encoding/gob"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"internal/timer"
	"google.golang.org/protobuf/proto"
	// "timer"
	// "github.com/tmck-code/pokesay-go/internal/timer"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func findFiles(dirpath string, ext string, skip []string) PokemonEntryMap {
	categories := &PokemonEntryMap{Pokemon: make(map[string]*PokemonEntryList)}
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

			for _, c := range pokemonCategories[3 : len(pokemonCategories)-1] {
				p := PokemonEntry{Name: fpath, Data: data}
				if val, ok := categories.Pokemon[c]; ok {
					val.Pokemon = append(val.Pokemon, &p)
				} else {
					categories.Pokemon[c] = &PokemonEntryList{Pokemon: []*PokemonEntry{&p}}
				}
				// fmt.Println("categories:", categories.Pokemon[c], c)

				categories.Pokemon[c].Pokemon = append(categories.Pokemon[c].Pokemon, &PokemonEntry{Name: fpath, Data: data})
			}
		}
		return err
	})
	check(err)

	return *categories
}

func createCategories(fpath string) []string {
	return strings.Split(fpath, "/")
}

func writeToFile(categories map[string][]PokemonEntry, fpath string) {
	ostream, err := os.Create(fpath)
	check(err)

	writer := bufio.NewWriter(ostream)
	enc := gob.NewEncoder(writer)
	enc.Encode(categories)
	writer.Flush()
	ostream.Close()
}

func writeProtobuf(categories PokemonEntryMap, fpath string) {
	fmt.Println("writing protobuf file")
	data, err := proto.Marshal(&categories)
	check(err)

	// printing out our raw protobuf object
	err = os.WriteFile("data.txt", data, 0644)
}

func readFromFile(fpath string) map[string][]PokemonEntry {
	istream, err := os.Open(fpath)
	check(err)

	reader := bufio.NewReader(istream)
	dec := gob.NewDecoder(reader)

	categories := make(map[string][]PokemonEntry)

	err = dec.Decode(&categories)
	check(err)
	istream.Close()

	return categories
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

func randInt(n int) int {
	if n <= 0 {
		log.Fatal("n <= 0! -_-")
	}
	return rand.Intn(n-1) + 1
}

func main() {
	args := parseArgs()
	fmt.Println("starting at", args.FromDir)

	t := timer.NewTimer()

	categories := findFiles(args.FromDir, ".cow", args.SkipDirs)

	// writeToFile(categories, args.ToFpath)
	writeProtobuf(categories, args.ToFpath)

	// categories := readFromFile(args.ToFpath)

	// categories := readFromEmbeddedString()
	t.Mark("writeProtobuf")

	total := 0
	for _, _ = range categories.Pokemon {
		total += 1
		// fmt.Println(k, len(v))
	}
	// fmt.Println("TOTAL:", total)
	rand.Seed(time.Now().UnixNano())
	randomCategory := randInt(total)
	t.Mark("randomCategory")

	idx := 0
	choice := ""
	for cat, _ := range categories.Pokemon {
		idx += 1
		if idx == randomCategory {
			choice = cat
		}
	}
	t.Mark("randomCategoryChoice")

	// fmt.Println(choice, categories[choice])
	fmt.Println("choosing pokemon", "cat:", randomCategory, "choice:", choice, " - ", len(categories.Pokemon[choice].Pokemon))
	randPokemon := randInt(len(categories.Pokemon[choice].Pokemon))
	t.Mark("randomPokemon")
	// fmt.Println("category choice", choice, "pokemon choice", randPokemon, "/", len(categories[choice]))
	pokemon := categories.Pokemon[choice].Pokemon[randPokemon]
	t.Mark("randomPokemonChoice")

	fmt.Printf("%s\n", pokemon.Data)
	t.Mark("print")
	t.StopTimer()
	t.PrintJson()
}
