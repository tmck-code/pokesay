package pokesay

import (
	"errors"
	"log"
	"math/rand"
	"sort"
	"time"

	"github.com/tmck-code/pokesay/src/pokedex"
)

var (
	Rand rand.Source = rand.NewSource(time.Now().UnixNano())
)

func RandomInt(n int) int {
	if n <= 0 {
		log.Fatal(errors.New("RandomInt arg must be >0"))
	}
	return rand.New(Rand).Intn(n)
}

func ChooseRandomCategory(keys [][]string, categories pokedex.PokemonTrie) ([]string, []*pokedex.PokemonEntry) {
	choice := keys[RandomInt(len(keys)-1)]
	category, err := categories.GetCategory(choice)
	Check(err)
	return choice, category
}

func ListCategories(categories pokedex.PokemonTrie) []string {
	ukm := map[string]bool{}
	for _, v := range categories.Keys {
		for _, k := range v {
			ukm[k] = true
		}
	}
	keys := make([]string, len(ukm))
	i := 0
	for k := range ukm {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}
