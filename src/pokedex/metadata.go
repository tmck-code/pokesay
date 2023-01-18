package pokedex

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

func (m Metadata) WriteToFile(fpath string) {
	WriteStructToFile(m, fpath)
}
