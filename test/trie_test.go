package test

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
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
					},"data":null}
				},"data":null}
			},"data":null
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
	t.Insert([]string{"p", "g2", "r"}, pokedex.NewEntry(2, "bulbasaur"))

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
            "data": null
          },
          "g2": {
            "children": {
              "r": {
                "children": {},
                "data": [
                  {
                    "value": "bulbasaur",
                    "index": 2
                  }
                ]
              }
            },
            "data": null
          }
        },
        "data": null
      }
    },
    "data": null
  },
  "len": 3,
  "keys": [
    [
      "p",
      "g1",
      "r"
    ],
    [
      "p",
      "g2",
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
	defer os.Remove("test.txt")

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
	Assert(expected, t.KeyPaths, t, test)

	expected = [][]string{
		{"big", "g1", "o"},
		{"big", "g1"},
	}
	result, err := t.FindKeyPaths("big")
	pokesay.Check(err)
	Assert(expected, result, result, test)
}

func TestFindByKeyPath(test *testing.T) {
	t := pokedex.NewTrie()

	t.Insert([]string{"small", "g1", "r"}, pokedex.NewEntry(0, "pikachu"))
	t.Insert([]string{"small", "g1", "o"}, pokedex.NewEntry(1, "bulbasaur"))
	t.Insert([]string{"medium", "g1", "o"}, pokedex.NewEntry(2, "bulbasaur"))
	t.Insert([]string{"big", "g1", "o"}, pokedex.NewEntry(3, "bulbasaur"))
	t.Insert([]string{"big", "g1"}, pokedex.NewEntry(4, "charmander"))

	result, err := t.FindByKeyPath([]string{"small", "g1"})
	pokesay.Check(err)

	expected := []*pokedex.Entry{{Index: 0, Value: "pikachu"}, {Index: 1, Value: "bulbasaur"}}

	for i := range result {
		Assert(expected[i], result[i], result[i], test)
	}
}

func TestWriteToFile(test *testing.T) {
	defer os.Remove("test.txt")

	t := pokedex.NewTrie()
	t.Insert([]string{"p", "g1", "r"}, pokedex.NewEntry(0, "pikachu"))
	t.Insert([]string{"p", "g1", "r"}, pokedex.NewEntry(1, "bulbasaur"))

	t.WriteToFile("test.txt")

	data, err := os.ReadFile("test.txt")
	pokesay.Check(err)

	d := &pokedex.Trie{}

	err = gob.NewDecoder(bytes.NewBuffer(data)).Decode(&d)
	pokesay.Check(err)

	Assert(t.ToString(), d.ToString(), d, test)
}
