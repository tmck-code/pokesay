package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/tmck-code/pokesay-go/src/pokedex"
	"github.com/tmck-code/pokesay-go/src/bin"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

var (
	DEBUG bool = false
)

type CowBuildArgs struct {
	FromDir  string
	ToDir    string
	SkipDirs []string
	Debug    bool
}

func parseArgs() CowBuildArgs {
	fromDir := flag.String("from", ".", "from dir")
	toDir := flag.String("to", ".", "to dir")
	skipDirs := flag.String("skip", "'[\"resources\"]'", "JSON array of dir patterns to skip converting")
	debug := flag.Bool("debug", DEBUG, "show debug logs")

	flag.Parse()

	DEBUG = *debug

	args := CowBuildArgs{FromDir: *fromDir, ToDir: *toDir}
	json.Unmarshal([]byte(*skipDirs), &args.SkipDirs)

	return args
}

func main() {
	args := parseArgs()

	fpaths := pokedex.FindFiles(args.FromDir, ".png", args.SkipDirs)

	fmt.Println("Converting PNGs -> cowfiles")
	pbar := bin.NewProgressBar(len(fpaths))
	for _, f := range fpaths {
		fmt.Println("~ Converting", f, "from", args.FromDir, "->", args.ToDir, 2)
		pokedex.ConvertPngToCow(args.FromDir, f, args.ToDir, 2)
		pbar.Add(1)
	}
	fmt.Println("Finished converting", len(fpaths), "pokesprite PNGs", "-> cowfiles")
}
