package pokedex

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"embed"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
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

func ReadStructFromBytes[T any](data []byte) T {
	var d T
	err := gob.NewDecoder(bytes.NewBuffer(data)).Decode(&d)
	Check(err)

	return d
}

func WriteStructToFile(obj interface{}, fpath string) {
	ostream, err := os.Create(fpath)
	Check(err)

	writer := bufio.NewWriter(ostream)
	gob.NewEncoder(writer).Encode(obj)

	writer.Flush()
	ostream.Close()
}

func StructToJSON(obj interface{}, indentation ...int) string {
	if len(indentation) == 1 {
		json, err := json.MarshalIndent(obj, "", strings.Repeat(" ", indentation[0]))
		Check(err)
		return string(json)
	} else {
		json, err := json.Marshal(obj)
		Check(err)
		return string(json)
	}
}

func WriteBytesToFile(data []byte, fpath string, compress bool) {
	ostream, err := os.Create(fpath)
	Check(err)

	writer := bufio.NewWriter(ostream)
	if compress {
		writer.WriteString(string(Compress(data)))
	} else {
		writer.WriteString(string(data))
	}
	writer.Flush()
	ostream.Close()
}

func WriteIntToFile(n int, fpath string) {
	WriteBytesToFile([]byte(strconv.Itoa(n)), fpath, false)
}

func ReadIntFromBytes(bs []byte) int {
	total, err := strconv.Atoi(string(bs))
	Check(err)

	return total
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

type Metadata struct {
	Data     []byte
	Index    int
	Metadata PokemonMetadata
}

func (m Metadata) WriteToFile(fpath string) {
	WriteStructToFile(m, fpath)
}

func CreateNameMetadata(idx int, key string, name PokemonName, rootDir string, fpaths []string) *PokemonMetadata {
	entryCategories := make(map[int][][]string, 0)

	for _, fpath := range fpaths {
		basename := strings.TrimPrefix(fpath, rootDir)
		if strings.Contains(basename, strings.ToLower(name.English)) {
			data, err := os.ReadFile(fpath)
			Check(err)
			cats := createCategories(strings.TrimPrefix(fpath, rootDir), data)
			entryCategories[idx] = append(entryCategories[idx], cats)
		}
	}
	return NewMetadata(
		name.English,
		name.Japanese,
		name.JapanesePhonetic,
		entryCategories,
	)
}

func CreateCategoryStruct(rootDir string, metadata []PokemonMetadata, debug bool) Trie {
	categories := NewTrie()
	for _, m := range metadata {
		for i, category := range m.Entries {
			categories.Insert(
				category.Categories,
				NewEntry(i, m.Name),
			)
		}
	}
	return *categories
}

func createName(fpath string) string {
	parts := strings.Split(fpath, "/")
	return strings.Split(parts[len(parts)-1], ".")[0]
}

func createCategories(fpath string, data []byte) []string {
	parts := strings.Split(fpath, "/")
	height := sizeCategory(len(strings.Split(string(data), "\n")))

	return append([]string{height}, parts[0:len(parts)-1]...)
}

func sizeCategory(height int) string {
	if height <= 13 {
		return "small"
	} else if height <= 19 {
		return "medium"
	}
	return "big"
}

func ReadPokemonCow(embeddedData embed.FS, fpath string) []byte {
	d, err := embeddedData.ReadFile(fpath)
	Check(err)

	return Decompress(d)
}
