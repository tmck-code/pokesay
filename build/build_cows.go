package main

import (
	"bytes"
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

var SENTINEL string = "__END__"
var N_FILES int

func findFiles(dirpath string, ext string) []string {
	fpaths := []string{}
	err := filepath.Walk(dirpath, func(path string, f os.FileInfo, err error) error {
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
	cmd := exec.Command("/usr/local/bin/img2xterm", sourceFpath)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &stderr

	err := cmd.Run()

	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		log.Fatal(err)
	}
	return out.Bytes()
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
	// fmt.Printf("%s -> %s\n", sourceFpath, destFpath)
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

	// var fpath string

	fpaths := findFiles(args.FromDir, ".png")
	pbar := newProgressBar(len(fpaths))

	fpathChan := make(chan string, len(fpaths))
	fmt.Println(len(fpaths), len(fpathChan))

	go func() {
		for _, f := range fpaths {
			fpathChan <- f
		}
		fpathChan <- SENTINEL
	}()

	var wg sync.WaitGroup
	var fpath string

	for true {
		fpath = <-fpathChan
		if fpath == SENTINEL {
			break
		}
		go convertPngToCow(args.FromDir, fpath, args.ToDir, &wg, &pbar)
		wg.Add(1)
	}
	wg.Wait()
	fmt.Println("finished converting", N_FILES, "pokemon .pngs -> cowfiles!")
}
