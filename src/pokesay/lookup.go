package pokesay

import (
	"embed"
	"fmt"
	"io/fs"
	"math/rand"
	"path"
	"strconv"
	"strings"
	"time"
	"log"

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

func GatherMapKeys(m map[string][]int) []string {
	keys := make([]string, 0)
	for k, _ := range m {
		keys = append(keys, k)
	}
	return keys
}

func ListNames(names map[string][]int) []string {
	return GatherMapKeys(names)
}

func ChooseByName(names map[string][]int, nameToken string, metadataFiles embed.FS, metadataRootDir string) (pokedex.PokemonMetadata, pokedex.PokemonEntryMapping) {
	match := names[nameToken]
	if len(match) == 0 {
		log.Fatal(fmt.Sprintf("cannot find pokemon by name '%s'", nameToken))
	}
	nameChoice := match[RandomInt(len(match))]

	metadata := pokedex.ReadMetadataFromEmbedded(
		metadataFiles,
		pokedex.MetadataFpath(metadataRootDir, nameChoice),
	)
	choice := RandomInt(len(metadata.Entries))
	return metadata, metadata.Entries[choice]
}

func ChooseByRandomIndex(totalInBytes []byte) (int, int) {
	total := pokedex.ReadIntFromBytes(totalInBytes)
	return total, RandomInt(total)
}
