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

func ChooseByCategory(categoryKey string, categories pokedex.Trie) (*pokedex.Entry, []string) {
	matches, err := categories.FindKeyPaths(categoryKey)
	Check(err)

	keyPath := matches[RandomInt(len(matches)-1)]
	category, err := categories.FindByKeyPath(keyPath)
	Check(err)

	choice := category[RandomInt(len(category))]

	return choice, keyPath
}

func ListCategories(categories pokedex.Trie) []string {
	ukm := map[string]bool{}
	for _, v := range categories.KeyPaths {
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

func ChooseByName(name string, categories pokedex.Trie) *pokedex.PokemonMatch {
	matches, err := categories.Find(name)
	Check(err)
	return matches[RandomInt(len(matches))]
}
