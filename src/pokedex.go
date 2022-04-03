package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
	"github.com/tmck-code/pokesay-go/src/pokedex"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func convertFiles(fpaths []string, pbar *progressbar.ProgressBar) (pokedex.PokemonTrie, []pokedex.Metadata) {
	categories := pokedex.NewTrie()
	metadata := []pokedex.Metadata{}
	for idx, fpath := range fpaths {
		data, err := os.ReadFile(fpath)
		check(err)

		categories.Insert(
			createCategories(fpath),
			pokedex.NewPokemonEntry(idx, createName(fpath)),
		)
		metadata = append(metadata, pokedex.Metadata{Data: data, Index: idx})
	}
	return *categories, metadata
}

func createName(fpath string) string {
	parts := strings.Split(fpath, "/")
	return strings.Split(parts[len(parts)-1], ".")[0]
}

func createCategories(fpath string) []string {
	parts := strings.Split(fpath, "/")
	return append([]string{"pokemon"}, parts[3:len(parts)-1]...)
}

type CowBuildArgs struct {
	FromDir    string
	ToDir      string
	SkipDirs   []string
	DebugTimer bool
}

func newProgressBar(max int) progressbar.ProgressBar {
	return *progressbar.NewOptions(
		max,
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionSetWidth(40),
		progressbar.OptionThrottle(100*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionOnCompletion(func() { fmt.Fprint(os.Stderr, "\n") }),
		progressbar.OptionSetTheme(progressbar.Theme{Saucer: "█", SaucerPadding: "░", BarStart: "╢", BarEnd: "╟"}),
	)
}

func parseArgs() CowBuildArgs {
	fromDir := flag.String("from", ".", "from dir")
	toDir := flag.String("to", ".", "to fpath")
	skipDirs := flag.String("skip", "'[\"resources\"]'", "JSON array of dir patterns to skip converting")
	debugTimer := flag.Bool("debugTimer", false, "show a debug timer")

	flag.Parse()

	args := CowBuildArgs{FromDir: *fromDir, ToDir: *toDir, DebugTimer: *debugTimer}
	json.Unmarshal([]byte(*skipDirs), &args.SkipDirs)

	return args
}

// func main() {
// 	args := parseArgs()
// 	t := timer.NewTimer()
// 	fmt.Println("starting at", args.FromDir)

// 	// categories is a PokemonTrie struct that will be written to a file using encoding/gob
// 	// metadata is a list of pokemon data and an index to use when writing them to a file
// 	// - this index matches a corresponding one in the categories struct
// 	// - these files are embedded into the build binary using go:embed and then loaded at runtime
// 	categories, metadata := convertFiles(
// 		FindFiles(args.FromDir, ".cow", args.SkipDirs),
// 	)
// 	t.Mark("CreateEntriesFromFiles")

// 	WriteStructToFile(categories, args.ToFpath)
// 	t.Mark("WriteCategoriesToFile")

// 	for _, m := range metadata {
// 		WriteCompressedToFile(m.Data, EntryFpath(m.Index))
// 	}
// 	t.Mark("WriteDataToFiles")

// 	if args.DebugTimer {
// 		t.Stop()
// 		t.PrintJson()
// 	}
// }

func main() {
	args := parseArgs()
	fmt.Println("starting at", args.FromDir)

	fpaths := pokedex.FindFiles(args.FromDir, ".png", args.SkipDirs)
	fpathChan := make(chan string, len(fpaths))

	go func() {
		for _, f := range fpaths {
			fpathChan <- f
		}
	}()

	var wg sync.WaitGroup
	pbar := newProgressBar(len(fpaths))

	for i := 0; i < len(fpaths); i++ {
		go pokedex.ConvertPngToCow(args.FromDir, <-fpathChan, args.ToDir, 2, &wg, &pbar)
		wg.Add(1)
	}
	wg.Wait()
	fmt.Println("Finished converting", len(fpaths), "pokesprite -> cowfiles")

	pbar = newProgressBar(len(fpaths))
	categories, metadata := convertFiles(fpaths, &pbar)
	fmt.Println("Writing categories to build/gob")
	pokedex.WriteStructToFile(categories, "build/gob")

	fmt.Println("Writing data files to build/*cow")
	for _, m := range metadata {
		pokedex.WriteCompressedToFile(m.Data, pokedex.EntryFpath(m.Index))
	}
}
