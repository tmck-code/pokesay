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

func NewMetadataFromBytes(data []byte) PokemonMetadata {
	return ReadStructFromBytes[PokemonMetadata](data)
}

func NewMetadataFromGOBData(data embed.FS, index int) PokemonMetadata {
	m, err := data.ReadFile(MetadataFpath("build/assets/metadata", index))
	Check(err)
	return NewMetadataFromBytes(m)
}
