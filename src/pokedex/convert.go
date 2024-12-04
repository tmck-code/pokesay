package pokedex

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	failures     []string
	COLOUR_RESET string = fmt.Sprintf("%s[%dm\n", "\x1b", 39)
)

func FindFiles(dirpath string, ext string, skip []string) []string {
	fpaths := []string{}
	err := filepath.Walk(dirpath, func(path string, f os.FileInfo, err error) error {
		for _, s := range skip {
			if strings.Contains(path, s) {
				return err
			}
		}
		if !f.IsDir() && filepath.Ext(f.Name()) == ext {
			fpaths = append(fpaths, path)
		}
		return err
	})
	Check(err)
	return fpaths
}

// img2xterm converts an image to a cowfile, returning the result as a byte slice
func img2xterm(sourceFpath string) ([]byte, error) {
	return exec.Command("bash", "-c", fmt.Sprintf("/usr/local/bin/img2xterm %s", sourceFpath)).Output()
}

// autoCrop trims the whitespace from the edges of an image, in place
func autoCrop(sourceFpath string) {
	destFpath := fmt.Sprintf("/tmp/%s", filepath.Base(sourceFpath))
	_, err := exec.Command(
		"bash", "-c", fmt.Sprintf("/usr/bin/convert %s -trim +repage %s", sourceFpath, destFpath),
	).Output()
	Check(err)

	os.Rename(destFpath, sourceFpath)
}

func countLineLeftPadding(line string) int {
	count := 0
	for _, ch := range line {
		if ch == ' ' {
			count += 1
		} else {
			break
		}
	}
	return count
}

func countCowfileLeftPadding(cowfile []byte) int {
	lines := strings.Split(string(cowfile), "\n")

	paddings := make([]int, 0)
	for _, line := range lines {
		paddings = append(paddings, countLineLeftPadding(line))
	}

	minPadding := 100 // TODO: a better way?
	for _, padding := range paddings {
		if padding > 0 && padding < minPadding {
			minPadding = padding
		}
	}
	return minPadding
}

func stripEmptyLines(cowfile []string) []string {
	converted := make([]string, 0)

	for _, line := range cowfile {
		if len(line) == 0 {
			continue
		}
		onlySpaces := true
		for _, ch := range line {
			if ch != ' ' {
				onlySpaces = false
				break
			}
		}
		if !onlySpaces {
			converted = append(converted, line)
		}
	}
	return converted
}

func padLeft(cowfile []byte, n int) []string {
	converted := make([]string, 0)
	lines := strings.Split(string(cowfile), "\n")

	for _, line := range lines {
		converted = append(converted, strings.Repeat(" ", n)+line)
	}
	return converted
}

func ConvertPngToCow(sourceDirpath string, sourceFpath string, destDirpath string, extraPadding int) {
	destDir := filepath.Join(
		destDirpath,
		// strip the root "source dirpath" from the source path
		// e.g. fpath: /a/b/c.txt sourceDir: /a/ -> b/c.txt
		filepath.Dir(strings.ReplaceAll(sourceFpath, sourceDirpath, "")),
	)
	// Ensure that the destination dir exists
	os.MkdirAll(destDir, 0755)

	// Trim the whitespace from the edges of the images. This helps with the conversion
	autoCrop(sourceFpath)
	// Some conversions are failing with something about colour channels
	// Can't be bothered resolving atm, so just skip past any failed conversions
	converted, _ := img2xterm(sourceFpath)

	if len(converted) == 0 {
		failures = append(failures, sourceFpath)
		return
	}

	destFpath := filepath.Join(destDir, strings.ReplaceAll(filepath.Base(sourceFpath), ".png", ".cow"))
	ostream, err := os.Create(destFpath)
	Check(err)
	defer ostream.Close()
	writer := bufio.NewWriter(ostream)

	final := stripEmptyLines(padLeft(converted, extraPadding))

	// Join all of the lines back together, add colour reset sequence at the end
	_, err = writer.WriteString(strings.Join(final, "\n") + COLOUR_RESET)
	Check(err)

	writer.Flush()
}
