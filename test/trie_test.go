package test

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/tmck-code/pokesay/src/pokedex"
	"github.com/tmck-code/pokesay/src/pokesay"
)

func TestNewEntry(test *testing.T) {
	p := pokedex.NewEntry(1, "yo")
	Assert(1, p.Index, p, test)
	Assert("yo", p.Value, p, test)
}

func TestTrieToString(test *testing.T) {
	t := pokedex.NewTrie()
	t.Insert([]string{"p", "g1", "r"}, pokedex.NewEntry(0, "pikachu"))
	t.Insert([]string{"p", "g1", "r"}, pokedex.NewEntry(1, "bulbasaur"))

	expected := FlattenJSON(`{
		"root":{
			"children":{
				"p":{"children":{
					"g1":{"children":{
						"r":{"children":{},
							"data":[{"value":"pikachu","index":0},{"value":"bulbasaur","index":1}]
						}
					},"data":[]}
				},"data":[]}
			},
			"data":null
		},
		"len":2,
		"keys":[["p","g1","r"]]
	}`)
	result := t.ToString()

	Assert(expected, string(result), string(result), test)
}

func TestTrieToStringIndented(test *testing.T) {
	t := pokedex.NewTrie()
	t.Insert([]string{"p", "g1", "r"}, pokedex.NewEntry(0, "pikachu"))
	t.Insert([]string{"p", "g1", "r"}, pokedex.NewEntry(1, "bulbasaur"))

	expected := `{
  "root": {
    "children": {
      "p": {
        "children": {
          "g1": {
            "children": {
              "r": {
                "children": {},
                "data": [
                  {
                    "value": "pikachu",
                    "index": 0
                  },
                  {
                    "value": "bulbasaur",
                    "index": 1
                  }
                ]
              }
            },
            "data": []
          }
        },
        "data": []
      }
    },
    "data": null
  },
  "len": 2,
  "keys": [
    [
      "p",
      "g1",
      "r"
    ]
  ]
}`
	result := t.ToString(2)

	Assert(expected, string(result), string(result), test)
}

func TestTrieInsert(test *testing.T) {
	t := pokedex.NewTrie()
	t.Insert([]string{"p", "g1", "r"}, pokedex.NewEntry(0, "pikachu"))
	t.Insert([]string{"p", "g1", "r"}, pokedex.NewEntry(1, "bulbasaur"))

	Assert(
		&pokedex.Entry{Value: "pikachu", Index: 0},
		t.Root.Children["p"].Children["g1"].Children["r"].Data[0],
		t,
		test,
	)
	Assert(
		&pokedex.Entry{Value: "bulbasaur", Index: 1},
		t.Root.Children["p"].Children["g1"].Children["r"].Data[1],
		t,
		test,
	)
}

func TestTrieFind(test *testing.T) {
	t := pokedex.NewTrie()
	t.Insert([]string{"p", "g1", "r"}, pokedex.NewEntry(0, "pikachu"))
	t.Insert([]string{"p", "g1", "r"}, pokedex.NewEntry(1, "bulbasaur"))
	t.Insert([]string{"x", "g1", "l"}, pokedex.NewEntry(2, "pikachu-other"))

	expected := []pokedex.PokemonMatch{
		{
			Entry: &pokedex.Entry{Index: 0, Value: "pikachu"},
			Keys:  []string{"p", "g1", "r"},
		},
		{
			Entry: &pokedex.Entry{Index: 2, Value: "pikachu-other"},
			Keys:  []string{"x", "g1", "l"},
		},
	}

	results, err := t.Find("pikachu")
	pokesay.Check(err)

	sort.Slice(results, func(i, j int) bool {
		return strings.Compare(results[i].Entry.Value, results[j].Entry.Value) == -1
	})

	for i := range results {
		Assert(expected[i].Entry, results[i].Entry, results[i], test)
	}
}

func TestTrieFindByKeyPath(test *testing.T) {
	t := pokedex.NewTrie()
	t.Insert([]string{"p", "g1", "r"}, pokedex.NewEntry(0, "pikachu"))
	t.Insert([]string{"p", "g1", "r"}, pokedex.NewEntry(1, "bulbasaur"))

	result, err := t.FindByKeyPath([]string{"p", "g1"})
	pokesay.Check(err)

	Assert(2, len(result), result, test)
	Assert(
		"[{Index: 0, Value: pikachu} {Index: 1, Value: bulbasaur}]",
		fmt.Sprintf("%s", result),
		result, test,
	)
}

func TestFindKeyPaths(test *testing.T) {
	t := pokedex.NewTrie()
	t.Insert([]string{"small", "g1", "r"}, pokedex.NewEntry(0, "pikachu"))
	t.Insert([]string{"small", "g1", "o"}, pokedex.NewEntry(1, "bulbasaur"))
	t.Insert([]string{"medium", "g1", "o"}, pokedex.NewEntry(2, "bulbasaur"))
	t.Insert([]string{"big", "g1", "o"}, pokedex.NewEntry(3, "bulbasaur"))
	t.Insert([]string{"big", "g1"}, pokedex.NewEntry(4, "charmander"))

	expected := [][]string{
		{"small", "g1", "r"},
		{"small", "g1", "o"},
		{"medium", "g1", "o"},
		{"big", "g1", "o"},
		{"big", "g1"},
	}
	Assert(expected, t.Keys, t, test)

	expected = [][]string{
		{"big", "g1", "o"},
		{"big", "g1"},
	}
	result, err := t.FindKeyPaths("big")
	pokesay.Check(err)
	Assert(expected, result, result, test)
}
