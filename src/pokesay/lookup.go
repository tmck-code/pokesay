package pokesay

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"math/rand"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/tmck-code/pokesay/src/pokedex"
)

var (
	Rand rand.Source = rand.NewSource(time.Now().UnixNano())
)

func RandomInt(n int) int {
	if n <= 0 {
		return 0
	}
	return rand.New(Rand).Intn(n)
}

func ChooseByCategory(category string, categoryDir []fs.DirEntry, categoryFiles embed.FS, categoryRootDir string, metadataFiles embed.FS, metadataRootDir string) (pokedex.PokemonMetadata, pokedex.PokemonEntryMapping) {
	if len(categoryDir) == 0 {
		log.Fatalf("cannot find pokemon by category '%s'", category)
	}
	choice := categoryDir[RandomInt(len(categoryDir))]

	categoryMetadata, err := categoryFiles.ReadFile(
		pokedex.CategoryFpath(categoryRootDir, category, choice.Name()),
	)
	Check(err)

	parts := strings.Split(string(categoryMetadata), "/")

	metadata := pokedex.ReadMetadataFromEmbedded(
		metadataFiles,
		path.Join(metadataRootDir, fmt.Sprintf("%s.metadata", parts[0])),
	)

	entryIndex, err := strconv.Atoi(string(parts[1]))
	Check(err)

	return metadata, metadata.Entries[entryIndex]
}

func ListNames(names map[string][]int) []string {
	return pokedex.GatherMapKeys(names)
}

func ChooseByName(names map[string][]int, nameToken string, metadataFiles embed.FS, metadataRootDir string, category string) (pokedex.PokemonMetadata, pokedex.PokemonEntryMapping) {
	match := names[nameToken]
	if len(match) == 0 {
		log.Fatalf("cannot find pokemon by name '%s'", nameToken)
	}
	nameChoice := match[RandomInt(len(match))]

	metadata := pokedex.ReadMetadataFromEmbedded(
		metadataFiles,
		pokedex.MetadataFpath(metadataRootDir, nameChoice),
	)

	// pick a random entry
	if category == "" {
		choice := RandomInt(len(metadata.Entries))
		return metadata, metadata.Entries[choice]
	// try to filter by desired category
	} else {
		matching := make([]pokedex.PokemonEntryMapping, 0)
		for _, entry := range metadata.Entries {
			log.Printf("entry: %v", entry)
			for _, entryCategory := range entry.Categories {
				if entryCategory == category {
					matching = append(matching, entry)
				}
			}
		}
		// if the category is not found for this pokemon, return a random entry
		if len(matching) == 0 {
			return metadata, metadata.Entries[RandomInt(len(metadata.Entries))]
		} else {
			return metadata, matching[RandomInt(len(matching))]
		}
	}
}

func ChooseByRandomIndex(totalInBytes []byte) (int, int) {
	total := pokedex.ReadIntFromBytes(totalInBytes)
	return total, RandomInt(total)
}
