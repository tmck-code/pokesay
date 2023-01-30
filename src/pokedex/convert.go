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

func img2xterm(sourceFpath string) ([]byte, error) {
	return exec.Command("bash", "-c", fmt.Sprintf("/usr/local/bin/img2xterm %s", sourceFpath)).Output()
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

func stripPadding(cowfile []byte, n int) []string {
	converted := make([]string, 0)
	lines := strings.Split(string(cowfile), "\n")

	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		convertedLine := ""
		for i, ch := range line {
			if i >= n {
				convertedLine += string(ch)
			}
		}
		if len(convertedLine) > 0 {
			converted = append(converted, convertedLine)
		}
	}
	return converted
}

func CreateDestination(sourceDirpath string, sourceFpath string, destDirpath string) string {
	destDir := filepath.Join(
		destDirpath,
		// strip the root "source dirpath" from the source path
		// e.g. fpath: /a/b/c.txt sourceDir: /a/ -> b/c.txt
		filepath.Dir(strings.ReplaceAll(sourceFpath, sourceDirpath, "")),
	)
	// Ensure that the destination dir exists
	os.MkdirAll(destDir, 0755)

	destFpath := filepath.Join(destDir, strings.ReplaceAll(filepath.Base(sourceFpath), ".png", ".cow"))

	return destFpath
}

func ConvertPngToCow(destFpath string, sourceFpath string, extraPadding int) string {
	// Some conversions are failing with something about colour channels
	// Can't be bothered resolving atm, so just skip past any failed conversions
	converted, _ := img2xterm(sourceFpath)

	if len(converted) == 0 {
		return ""
	}
	final := stripPadding(converted, countCowfileLeftPadding(converted)-extraPadding)

	return strings.Join(final, "\n") + COLOUR_RESET
}

func WriteCowfile(data string, fpath string) {
	ostream, err := os.Create(fpath)
	defer ostream.Close()
	Check(err)

	writer := bufio.NewWriter(ostream)
	_, err = writer.WriteString(data)
	Check(err)

	writer.Flush()
}
