package pokedex

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"os"
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

func ReadMetadataFromBytes(data []byte) PokemonMetadata {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)

	d := &PokemonMetadata{}

	err := dec.Decode(&d)
	Check(err)

	return *d
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
