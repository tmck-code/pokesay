package pokesay

import (
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
		return 0
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
	matches, err := categories.Find(name, true)
	Check(err)
	// for i, m := range matches {
	// 	fmt.Printf("name matches for %s: %v (%d)\n", name, m, i)
	// }
	choice := matches[RandomInt(len(matches))]

	// fmt.Printf("name choice for %s\n", choice)
	return choice
}

func ChooseByRandomIndex(totalInBytes []byte) (int, int) {
	total := pokedex.ReadIntFromBytes(totalInBytes)
	return total, RandomInt(total)
}
