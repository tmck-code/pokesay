package pokedex

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	Failures     []string
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
	return exec.Command("bash", "-c", fmt.Sprintf("/usr/local/bin/img2xterm %s 2>&1", sourceFpath)).Output()
}

// autoCrop trims the whitespace from the edges of an image, in place
func autoCrop(sourceFpath string, tmpDirpath string) (string, error) {
	// Generate a random suffix to avoid conflicts when running in parallel
	randomBytes := make([]byte, 8)
	rand.Read(randomBytes)
	randomSuffix := hex.EncodeToString(randomBytes)

	destFpath := fmt.Sprintf("%s/%s-%s", tmpDirpath, randomSuffix, filepath.Base(sourceFpath))
	// fmt.Println("Auto-cropping", sourceFpath, "->", destFpath)
	output, err := exec.Command(
		"bash", "-c", fmt.Sprintf("/usr/bin/convert %s -trim +repage %s 2>&1", sourceFpath, destFpath),
	).Output()
	if err != nil {
		return "", fmt.Errorf("auto-crop failed for %s: (%v) - %s", sourceFpath, err, strings.Trim(string(output), "\n"))
	}

	return destFpath, nil
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

func ConvertPngToCow(sourceDirpath string, sourceFpath string, tmpDirpath string, destDirpath string, extraPadding int) (string, error) {

	// Trim the whitespace from the edges of the images. This helps with the conversion
	tmpFpath, err := autoCrop(sourceFpath, tmpDirpath)
	if err != nil {
		// If autoCrop fails, still try to convert the original image
		Failures = append(Failures, string(err.Error()))
		return "", err
	}

	// Some conversions are failing with something about colour channels
	output, err := img2xterm(tmpFpath)
	if err != nil {
		failureMsg := fmt.Sprintf("failed to convert %s: (%v) - %s", tmpFpath, err, strings.Trim(string(output), "\n"))
		Failures = append(Failures, failureMsg)
		return "", errors.New(failureMsg)
	}

	if len(output) == 0 {
		failureMsg := fmt.Sprintf("failed to convert %s: no output", tmpFpath)
		Failures = append(Failures, failureMsg)
		return "", errors.New(failureMsg)
	}
	final := stripEmptyLines(padLeft(output, extraPadding))
	return strings.Join(final, "\n") + COLOUR_RESET, nil
}

func WriteToCowfile(data string, destDirpath string, destFpath string) {
	// Ensure that the destination dir exists
	err := os.MkdirAll(destDirpath, 0755)
	if err != nil {
		fmt.Printf("Failed to create directory %s: %v\n", destDirpath, err)
		Check(err)
	}

	ostream, err := os.Create(destFpath)
	if err != nil {
		fmt.Printf("Failed to create file %s: %v\n", destFpath, err)
		Check(err)
	}
	defer ostream.Close()
	writer := bufio.NewWriter(ostream)

	_, err = writer.WriteString(data)
	Check(err)

	err = writer.Flush()
	Check(err)
}
