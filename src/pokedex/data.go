package pokedex

import (
	"os"
	"log"
	"bufio"
	"encoding/json"
	"fmt"
	"strings"
)

// "name": {
//   "eng": "Caterpie",
//   "chs": "绿毛虫",
//   "jpn": "キャタピー",
//   "jpn_ro": "Caterpie"
// },
type DataEntryName struct {
	Eng string `json:"eng"`
	Chs string `json:"chs"`
	Jpn string `json:"jpn"`
	Jpn_ro string `json:"jpn_ro"`
}
// "slug":{
//   "eng":"bulbasaur","
//   "jpn":"fushigidane",
//   "jpn_ro":"fushigidane"
// }
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
	JapaneseRomaji string
}


func NewPokemonName(entry DataEntry) *PokemonName {
	return &PokemonName{
		English: entry.Name.Eng,
		Japanese: entry.Name.Jpn,
		JapaneseRomaji: entry.Slug.Jpn_ro,
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

