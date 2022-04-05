package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
	"strconv"

	"github.com/schollz/progressbar/v3"
	"github.com/tmck-code/pokesay-go/src/pokedex"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

type CowBuildArgs struct {
	FromDir    string
	ToDir      string
	SkipDirs   []string
	DebugTimer bool
	ToCategoryFpath string
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
	toDir := flag.String("to", ".", "to dir")
	toCategoryFpath := flag.String("toCategoryFpath", "build/pokedex.gob", "to fpath")
	skipDirs := flag.String("skip", "'[\"resources\"]'", "JSON array of dir patterns to skip converting")
	debugTimer := flag.Bool("debugTimer", false, "show a debug timer")

	flag.Parse()

	args := CowBuildArgs{FromDir: *fromDir, ToDir: *toDir, ToCategoryFpath: *toCategoryFpath, DebugTimer: *debugTimer}
	json.Unmarshal([]byte(*skipDirs), &args.SkipDirs)

	return args
}

func main() {
	args := parseArgs()
	fmt.Println("starting at", args.FromDir)

	fpaths := pokedex.FindFiles(args.FromDir, ".png", args.SkipDirs)

	pbar := newProgressBar(len(fpaths))
	for _, f := range fpaths {
		pokedex.ConvertPngToCow(args.FromDir, f, args.ToDir, 2, &pbar)
	}
	fmt.Println("Finished converting", len(fpaths), "pokesprite -> cowfiles")

	fmt.Println("starting at", args.ToDir)
	fpaths = pokedex.FindFiles(args.ToDir, ".cow", args.SkipDirs)

	// categories is a PokemonTrie struct that will be written to a file using encoding/gob
	// metadata is a list of pokemon data and an index to use when writing them to a file
	// - this index matches a corresponding one in the categories struct
	// - these files are embedded into the build binary using go:embed and then loaded at runtime
	categories, metadata := pokedex.CreateMetadata(fpaths)

	fmt.Println("writing categories to", args.ToCategoryFpath)
	pokedex.WriteStructToFile(categories, args.ToCategoryFpath)

	for _, m := range metadata {
		pokedex.WriteBytesToFile(m.Data, pokedex.EntryFpath(m.Index), true)
		pokedex.WriteStructToFile(m.Metadata, pokedex.MetadataFpath(m.Index))
	}
	pokedex.WriteBytesToFile([]byte(strconv.Itoa(len(metadata))), "build/total.txt", false)
}
