package pokedex

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
)

func Check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func EntryFpath(subdir string, idx int) string {
	return path.Join(subdir, fmt.Sprintf("%d.cow", idx))
}

func MetadataFpath(subdir string, idx int) string {
	return path.Join(subdir, fmt.Sprintf("%d.metadata", idx))
}

func Compress(data []byte) []byte {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)

	_, err := gz.Write(data)
	Check(err)

	err = gz.Close()
	Check(err)

	return b.Bytes()
}

func Decompress(data []byte) []byte {
	buf := bytes.NewBuffer(data)

	reader, err := gzip.NewReader(buf)
	Check(err)

	var resB bytes.Buffer

	_, err = resB.ReadFrom(reader)
	Check(err)

	return resB.Bytes()
}

func CreateMetadata(rootDir string, fpaths []string, pokemonNames map[string]PokemonName, debug bool) []Metadata {
	metadata := []Metadata{}
	for i, fpath := range fpaths {
		data, err := os.ReadFile(fpath)
		Check(err)

		cats := createCategories(strings.TrimPrefix(fpath, rootDir), data)
		name := createName(fpath)

		v, _ := pokemonNames[strings.Split(name, "-")[0]]

		fmt.Println("name:", name, "found:", v)

		metadata = append(
			metadata,
			Metadata{
				data,
				i,
				PokemonMetadata{
					Name:             name,
					JapaneseName:     v.Japanese,
					JapanesePhonetic: v.JapanesePhonetic,
					Categories:       strings.Join(cats, "/"),
				},
			},
		)
	}
	return metadata
}

func CreateCategoryStruct(rootDir string, fpaths []string, debug bool) PokemonTrie {
	categories := NewTrie()
	for i, fpath := range fpaths {
		data, err := os.ReadFile(fpath)
		Check(err)

		cats := createCategories(strings.TrimPrefix(fpath, rootDir), data)
		name := createName(fpath)

		fmt.Println("name:", name)

		categories.Insert(
			cats,
			NewPokemonEntry(i, name),
		)
	}
	return *categories
}

func createName(fpath string) string {
	parts := strings.Split(fpath, "/")
	return strings.Split(parts[len(parts)-1], ".")[0]
}

func SizeCategory(height int) string {
	if height <= 13 {
		return "small"
	} else if height <= 19 {
		return "medium"
	}
	return "big"
}

func createCategories(fpath string, data []byte) []string {
	parts := strings.Split(fpath, "/")
	height := SizeCategory(len(strings.Split(string(data), "\n")))

	return append([]string{height}, parts[0:len(parts)-1]...)
}
