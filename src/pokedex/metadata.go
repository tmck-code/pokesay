package pokedex

import (
	"bytes"
	"encoding/gob"
)

type Metadata struct {
	Data     []byte
	Index    int
	Metadata PokemonMetadata
}

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
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)

	d := &PokemonMetadata{}

	err := dec.Decode(&d)
	Check(err)

	return *d
}

func (m Metadata) WriteToFile(fpath string) {
	WriteStructToFile(m, fpath)
}
