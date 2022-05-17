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
	check(err)
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
		converted = append(converted, convertedLine)
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

	// Some conversions are failing with something about colour channels
	// Can't be bothered resolving atm, so just skip past any failed conversions
	converted, _ := img2xterm(sourceFpath)

	if len(converted) == 0 {
		failures = append(failures, sourceFpath)
		return
	}

	destFpath := filepath.Join(destDir, strings.ReplaceAll(filepath.Base(sourceFpath), ".png", ".cow"))
	ostream, err := os.Create(destFpath)
	check(err)
	defer ostream.Close()
	writer := bufio.NewWriter(ostream)

	final := stripPadding(converted, countCowfileLeftPadding(converted)-extraPadding)

	// Join all of the lines back together, add colour reset sequence at the end
	_, err = writer.WriteString(strings.Join(final, "\n") + COLOUR_RESET)
	check(err)

	writer.Flush()
}

type Metadata struct {
	Data     []byte
	Index    int
	Metadata PokemonMetadata
}

func CreateMetadata(rootDir string, fpaths []string, debug bool) (PokemonTrie, []Metadata) {
	categories := NewTrie()
	metadata := []Metadata{}
	for i, fpath := range fpaths {
		data, err := os.ReadFile(fpath)
		check(err)

		cats := createCategories(strings.TrimPrefix(fpath, rootDir), data)
		name := createName(fpath)

		categories.Insert(
			cats,
			NewPokemonEntry(i, name),
		)
		metadata = append(metadata, Metadata{data, i, PokemonMetadata{Name: name, Categories: strings.Join(cats, "/")}})
	}
	return *categories, metadata
}

func createName(fpath string) string {
	parts := strings.Split(fpath, "/")
	return strings.Split(parts[len(parts)-1], ".")[0]
}

func SizeCategory(height int) string {
	if height <= 13 {
		return "small"
	} else if height <= 19 {
		return "medium"
	}
	return "big"
}

func createCategories(fpath string, data []byte) []string {
	parts := strings.Split(fpath, "/")
	height := SizeCategory(len(strings.Split(string(data), "\n")))

	return append([]string{height}, parts[0:len(parts)-1]...)
}

// Strips the leading "./" from a path e.g. "./cows/ -> cows/"
func NormaliseRelativeDir(dirPath string) string {
	return strings.TrimPrefix(dirPath, "./")
}
