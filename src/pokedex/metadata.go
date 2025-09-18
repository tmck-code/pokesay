package pokedex

import (
	"embed"
	"os"

	"github.com/tmck-code/pokesay/src/timer"
)

type PokemonEntryMapping struct {
	EntryIndex int
	Categories []string
}

type PokemonMetadata struct {
	Name             string
	JapaneseName     string
	JapanesePhonetic string
	Entries          []PokemonEntryMapping
}

func NewMetadata(name string, japaneseName string, japanesePhonetic string, entryMap map[int][][]string) *PokemonMetadata {

	entries := make([]PokemonEntryMapping, 0)

	for idx, categories := range entryMap {
		for _, category := range categories {
			entries = append(entries, PokemonEntryMapping{idx, category})
		}
	}

	return &PokemonMetadata{
		Name:             name,
		JapaneseName:     japaneseName,
		JapanesePhonetic: japanesePhonetic,
		Entries:          entries,
	}
}

func ReadMetadataFromBytes(data []byte) PokemonMetadata {
	return ReadStructFromBytes[PokemonMetadata](data)
}

func ReadMetadataFromFile(fpath string) PokemonMetadata {
	metadata, err := os.ReadFile(fpath)
	Check(err)
	timer.DebugTimer.Mark("read file")

	data := ReadMetadataFromBytes(metadata)
	timer.DebugTimer.Mark("read metadata")

	return data
}

func ReadMetadataFromEmbedded(embeddedData embed.FS, fpath string) PokemonMetadata {
	metadata, err := embeddedData.ReadFile(fpath)
	Check(err)
	timer.DebugTimer.Mark("read embedded file")

	data := ReadMetadataFromBytes(metadata)
	timer.DebugTimer.Mark("metadata from bytes")

	return data
}
