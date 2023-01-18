package pokedex

import (
	"os"
	"log"
	"bufio"
	"encoding/json"
	"fmt"
	"strings"
)

// {
//   "name": { "eng": "Charizard", "chs": "喷火龙", "jpn": "リザードン", "jpn_ro": "Lizardon" }
//   "slug": { "eng": "charizard",                  "jpn": "riza-don",   "jpn_ro": "lizardon" }
// }
// Out of all these names, we want the name.jpn, name.jpn_ro, slug.eng,
type DataEntryName struct {
	Eng string `json:"eng"`
	Jpn string `json:"jpn"`
	Jpn_ro string `json:"jpn_ro"`
}

type DataEntrySlug struct {
	Eng string `json:"eng"`
	Jpn string `json:"jpn"`
	Jpn_ro string `json:"jpn_ro"`
}

type DataEntry struct {
	Name DataEntryName `json:"name"`
	Slug DataEntrySlug `json:"slug"`
}

type PokemonName struct {
	English string
	Japanese string
	JapanesePhonetic string
}


func NewPokemonName(entry DataEntry) *PokemonName {
	return &PokemonName{
		English: entry.Name.Eng,
		Japanese: entry.Name.Jpn,
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
    // optionally, resize scanner's capacity for lines over 64K, see next example
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

