package pokedex

import (
	"embed"

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

func ReadMetadataFromEmbedded(embeddedData embed.FS, fpath string) PokemonMetadata {
	t := timer.NewTimer("ReadMetadataFromEmbedded")
	metadata, err := embeddedData.ReadFile(fpath)
	Check(err)
	t.Mark("read file")

	data := ReadMetadataFromBytes(metadata)
	t.Mark("read metadata")

	t.Stop()
	t.PrintJson()
	return data
}
