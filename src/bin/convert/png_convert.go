package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"

	"github.com/tmck-code/pokesay/src/bin"
	"github.com/tmck-code/pokesay/src/pokedex"
)

var (
	DEBUG bool = os.Getenv("DEBUG") != ""
)

type CowBuildArgs struct {
	FromDir  string
	TmpDir   string
	ToDir    string
	SkipDirs []string
	Padding  int
	Debug    bool
}

func parseArgs() CowBuildArgs {
	fromDir := flag.String("from", ".", "from dir")
	tmpDir := flag.String("tmpDir", "/tmp/convert/", "temporary directory for intermediate files")
	toDir := flag.String("to", ".", "to dir")
	skipDirs := flag.String("skip", "'[\"resources\"]'", "JSON array of dir patterns to skip converting")
	padding := flag.Int("padding", 2, "the number of spaces to pad from the left")

	flag.Parse()

	args := CowBuildArgs{FromDir: *fromDir, TmpDir: *tmpDir, ToDir: *toDir, Padding: *padding, Debug: DEBUG}
	json.Unmarshal([]byte(*skipDirs), &args.SkipDirs)

	if args.Debug {
		fmt.Printf("%+v\n", args)
	}
	return args
}

func worker(args CowBuildArgs, jobs <-chan string, pbar *progressbar.ProgressBar, dataSet map[string]struct{}, nDuplicates *int, nFailures *int, mu *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()

	for f := range jobs {
		data, err := pokedex.ConvertPngToCow(args.FromDir, f, args.TmpDir, args.ToDir, args.Padding)
		pbar.Add(1)

		if err != nil {
			mu.Lock()
			*nFailures++
			mu.Unlock()
			continue
		}

		// check if this cawfile is a duplicate of one that has already been written
		mu.Lock()
		if _, found := dataSet[data]; found {
			if args.Debug {
				fmt.Println("Skipping duplicate data for", f)
			}
			*nDuplicates++
			mu.Unlock()
			continue
		}
		dataSet[data] = struct{}{}
		mu.Unlock()

		destDirpath := filepath.Join(
			args.ToDir,
			// strip the root "source dirpath" from the source path
			// e.g. fpath: /a/b/c.txt sourceDir: /a/ -> b/c.txt
			filepath.Dir(strings.ReplaceAll(f, args.FromDir, "")),
		)
		destFpath := filepath.Join(destDirpath, strings.ReplaceAll(filepath.Base(f), ".png", ".cow"))

		pokedex.WriteToCowfile(data, destDirpath, destFpath)
	}
}

func main() {
	args := parseArgs()

	fpaths := pokedex.FindFiles(args.FromDir, ".png", args.SkipDirs)

	// Ensure that the destination dir exists
	os.MkdirAll(args.ToDir, 0755)

	fmt.Println("Converting PNGs -> cowfiles")
	pbar := bin.NewProgressBar(len(fpaths))

	dataSet := make(map[string]struct{})
	nDuplicates, nFailures := 0, 0

	nWorkers := runtime.NumCPU()
	var wg sync.WaitGroup
	var mu sync.Mutex

	// Create a channel to distribute work
	jobs := make(chan string, len(fpaths))

	// Send all file paths to the jobs channel
	for _, f := range fpaths {
		jobs <- f
	}
	close(jobs)

	// Start worker goroutines
	for w := 0; w < nWorkers; w++ {
		wg.Add(1)
		go worker(args, jobs, &pbar, dataSet, &nDuplicates, &nFailures, &mu, &wg)
	}

	// Wait for all workers to finish
	wg.Wait()
	fmt.Println("\nFinished converting", len(fpaths), "pokesprite PNGs -> cowfiles")
	fmt.Println("(skipped", nDuplicates, "duplicates and", nFailures, "failures)")

	// wait for progress bar to finish
	time.Sleep(100 * time.Millisecond)

	if args.Debug && len(pokedex.Failures) > 0 {
		fmt.Println("failures:")
		for _, f := range pokedex.Failures {
			fmt.Println(" -", f)
		}
	}
}
