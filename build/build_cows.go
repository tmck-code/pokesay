package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
)

type CowBuildArgs struct {
	FromDir  string
	ToDir    string
	SkipDirs []string
}

func parseArgs() CowBuildArgs {
	fromDir := flag.String("from", ".", "from dir")
	toDir := flag.String("to", ".", "to dir")
	skipDirs := flag.String("skip", "'[\"resources\"]'", "JSON array of dir patterns to skip converting")

	flag.Parse()

	args := CowBuildArgs{FromDir: *fromDir, ToDir: *toDir}
	json.Unmarshal([]byte(*skipDirs), &args.SkipDirs)

	return args
}

func findFiles(dirpath string, ext string, skip []string) []string {
	fpaths := []string{}
	err := filepath.Walk(dirpath, func(path string, f os.FileInfo, err error) error {
		for _, s := range(skip) {
			if strings.Contains(path, s) {
				return err
			}
		}
		if !f.IsDir() && filepath.Ext(f.Name()) == ext {
			fpaths = append(fpaths, path)
		}
		return err
	})
	if err != nil {
		fmt.Println("Fatal!")
		log.Fatal(err)
	}
	return fpaths
}

func img2xterm(sourceFpath string) []byte {
	out, err := exec.Command("bash", "-c", fmt.Sprintf("/usr/local/bin/img2xterm %s | grep \"\\S\"", sourceFpath)).Output()

	if err != nil {
		log.Fatal(err)
	}
	return out
}

func convertPngToCow(sourceDirpath string, sourceFpath string, destDirpath string, wg *sync.WaitGroup, pbar *progressbar.ProgressBar) {
	defer wg.Done()
	destDir := filepath.Join(
		destDirpath,
		// strip the root "source dirpath" from the source path
		// e.g. fpath: /a/b/c.txt sourceDir: /a/ -> b/c.txt
		filepath.Dir(strings.ReplaceAll(sourceFpath, sourceDirpath, "")),
	)
	// Ensure that the destination dir exists
	os.MkdirAll(destDir, 0755)
	time.Sleep(0)

	destFpath := filepath.Join(destDir, strings.ReplaceAll(filepath.Base(sourceFpath), ".png", ".cow"))

	err := os.WriteFile(destFpath, img2xterm(sourceFpath), 0644)
	time.Sleep(0)
	if err != nil {
		log.Fatal(err)
	}
	pbar.Add(1)
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

func main() {
	args := parseArgs()
	fmt.Println("starting at", args.FromDir)

	fpaths := findFiles(args.FromDir, ".png", args.SkipDirs)
	fpathChan := make(chan string, len(fpaths))

	go func() {
		for _, f := range fpaths {
			fpathChan <- f
		}
	}()

	var wg sync.WaitGroup
	pbar := newProgressBar(len(fpaths))

	for i := 0; i < len(fpaths) ; i++ {
		go convertPngToCow(args.FromDir, <-fpathChan, args.ToDir, &wg, &pbar)
		wg.Add(1)
	}
	wg.Wait()
}
