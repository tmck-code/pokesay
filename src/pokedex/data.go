package pokedex

import (
	"os"
	"log"
	"bufio"
	"encoding/json"
	"fmt"
)

type PokemonName struct {
	Eng string `json:"eng"`
	Chs string `json:"chs"`
	Jpn string `json:"jpn"`
	Jpn_ro string `json:"jpn_ro"`
}

// "name": {
//   "eng": "Caterpie",
//   "chs": "绿毛虫",
//   "jpn": "キャタピー",
//   "jpn_ro": "Caterpie"
// },
type DataEntry struct {
	Name PokemonName `json:"name"`
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
		entries[entry.Name.Eng] = entry.Name
    }
    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
	return entries
}

