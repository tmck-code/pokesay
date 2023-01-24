package pokedex

import "embed"

type PokemonMetadata struct {
	Categories       string
	Name             string
	JapaneseName     string
	JapanesePhonetic string
}

func NewMetadata(categories string, name string, japaneseName string, japanesePhonetic string) *PokemonMetadata {
	return &PokemonMetadata{
		Categories:       categories,
		Name:             name,
		JapaneseName:     japaneseName,
		JapanesePhonetic: japanesePhonetic,
	}
}

func ReadMetadataFromBytes(data []byte) PokemonMetadata {
	return ReadStructFromBytes[PokemonMetadata](data)
}

func ReadMetadataFromEmbedded(embeddedData embed.FS, fpath string) PokemonMetadata {
	metadata, err := embeddedData.ReadFile(fpath)
	Check(err)
	return ReadMetadataFromBytes(metadata)
}
