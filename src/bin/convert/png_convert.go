package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/tmck-code/pokesay/src/bin"
	"github.com/tmck-code/pokesay/src/pokedex"
)

var (
	DEBUG bool = false
)

type CowBuildArgs struct {
	FromDir  string
	ToDir    string
	SkipDirs []string
	Padding  int
	Debug    bool
}

func parseArgs() CowBuildArgs {
	fromDir := flag.String("from", ".", "from dir")
	toDir := flag.String("to", ".", "to dir")
	skipDirs := flag.String("skip", "'[\"resources\"]'", "JSON array of dir patterns to skip converting")
	padding := flag.Int("padding", 2, "the number of spaces to pad from the left")
	debug := flag.Bool("debug", DEBUG, "show debug logs")

	flag.Parse()

	DEBUG = *debug

	args := CowBuildArgs{FromDir: *fromDir, ToDir: *toDir, Padding: *padding}
	json.Unmarshal([]byte(*skipDirs), &args.SkipDirs)

	if DEBUG {
		fmt.Println("%+v", args)
	}
	return args
}

func main() {
	args := parseArgs()

	fpaths := pokedex.FindFiles(args.FromDir, ".png", args.SkipDirs)

	// Ensure that the destination dir exists
	os.MkdirAll(args.ToDir, 0755)

	fmt.Println("Converting PNGs -> cowfiles")
	pbar := bin.NewProgressBar(len(fpaths))
	for _, f := range fpaths {
		pokedex.ConvertPngToCow(args.FromDir, f, args.ToDir, args.Padding)
		pbar.Add(1)
	}
	fmt.Println("Finished converting", len(fpaths), "pokesprite PNGs", "-> cowfiles")
}
