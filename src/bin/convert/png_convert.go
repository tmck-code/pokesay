package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tmck-code/pokesay/src/bin"
	"github.com/tmck-code/pokesay/src/pokedex"
)

var (
	DEBUG bool = os.Getenv("DEBUG") != ""
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

	flag.Parse()

	args := CowBuildArgs{FromDir: *fromDir, ToDir: *toDir, Padding: *padding, Debug: DEBUG}
	json.Unmarshal([]byte(*skipDirs), &args.SkipDirs)

	if args.Debug {
		fmt.Printf("%+v\n", args)
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

	allData := make([]string, 0, len(fpaths))
	nDuplicates, nFailures := 0, 0

	for _, f := range fpaths {
		data, err := pokedex.ConvertPngToCow(args.FromDir, f, args.ToDir, args.Padding)
		if err != nil {
			nFailures++
			continue
		}

		// check if this cawfile is a duplicate of one that has already been written
		found := false
		for _, existingData := range allData {
			if existingData == data {
				found = true
				break
			}
		}
		if found {
			if args.Debug {
				fmt.Println("Skipping duplicate data for", f)
			}
			nDuplicates++
			pbar.Add(1)
			continue
		}
		allData = append(allData, data)

		destDirpath := filepath.Join(
			args.ToDir,
			// strip the root "source dirpath" from the source path
			// e.g. fpath: /a/b/c.txt sourceDir: /a/ -> b/c.txt
			filepath.Dir(strings.ReplaceAll(f, args.FromDir, "")),
		)
		destFpath := filepath.Join(destDirpath, strings.ReplaceAll(filepath.Base(f), ".png", ".cow"))

		pokedex.WriteToCowfile(data, destDirpath, destFpath)
		pbar.Add(1)
	}
	fmt.Println("\nFinished converting", len(fpaths), "pokesprite PNGs -> cowfiles")
	fmt.Println("(skipped", nDuplicates, "duplicates and", nFailures, "failures)")
}
