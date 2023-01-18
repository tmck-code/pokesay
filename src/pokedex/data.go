package pokedex

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

//	{
//	  "name": { "eng": "Charizard", "chs": "喷火龙", "jpn": "リザードン", "jpn_ro": "Lizardon" }
//	  "slug": { "eng": "charizard",                  "jpn": "riza-don",   "jpn_ro": "lizardon" }
//	}
//
// Out of all these names, we want the name.jpn, name.jpn_ro, slug.eng,
type DataEntry struct {
	Name struct {
		Eng    string `json:"eng"`
		Jpn    string `json:"jpn"`
		Jpn_ro string `json:"jpn_ro"`
	} `json:"name"`
	Slug struct {
		Eng    string `json:"eng"`
		Jpn    string `json:"jpn"`
		Jpn_ro string `json:"jpn_ro"`
	} `json:"slug"`
}

type PokemonName struct {
	English          string
	Japanese         string
	JapanesePhonetic string
}

func NewPokemonName(entry DataEntry) *PokemonName {
	return &PokemonName{
		English:          entry.Name.Eng,
		Japanese:         entry.Name.Jpn,
		JapanesePhonetic: entry.Slug.Jpn,
	}
}

func ReadNames(fpath string) map[string]PokemonName {
	istream, err := os.Open(fpath)
	if err != nil {
		log.Fatal(err)
	}
	defer istream.Close()

	entries := make(map[string]PokemonName)
	scanner := bufio.NewScanner(istream)
	for scanner.Scan() {
		line := scanner.Bytes()
		var entry DataEntry
		jsonErr := json.Unmarshal(line, &entry)
		if jsonErr != nil {
			fmt.Println(jsonErr)
		}
		entries[strings.ToLower(entry.Name.Eng)] = *NewPokemonName(entry)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return entries
}
