package paths

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const tilde = `~`

var startingPathDirs = []string{"~/.config/pathhelper/paths.d",
	"/etc/paths",
	"/etc/paths.d"}
var startingManPathDirs = []string{"~/.config/pathhelper/manpaths.d",
	"/etc/manpaths",
	"/etc/manpaths.d"}

var finalPaths []string
var finalManPaths []string

func init() {
	finalPaths = extractPaths(startingPathDirs)
	finalManPaths = extractPaths(startingManPathDirs)
}

func verify(path string) (err error) {

	return
}

func verifyPath(path string) (err error) {
	if _, err = os.Stat(path); err == nil {
		return
	} else if errors.Is(err, os.ErrNotExist) {
		return
	}

	return
}

func getPathsInDir(path string) (paths []string, err error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}

	for _, file := range files {
		if !file.IsDir() {
			err = verify(file.Name())
			if err != nil {
				continue
			}
			path := filepath.Join(path, file.Name())
			err = verifyPath(path)
			if err != nil {
				continue
			}
			paths = append(paths, path)
		}
	}

	return
}

// Paths get list of paths
func Paths() (paths []string) {
	return finalPaths
}

// ManPaths get list of man paths
func ManPaths() (paths []string) {
	return finalManPaths
}

// ZshFormatPath get zsh format path declaration
func ZshFormatPath() (output string) {
	return zshFormat("PATH")
}

// ZshFormatManPath get zsh format man path declaration
func ZshFormatManPath() (output string) {
	return zshFormat("MANPATH")
}

func zshFormat(envVar string) (output string) {
	output = fmt.Sprintf("export %s=\"%s\";", envVar, strings.Join(finalPaths, ":"))

	return
}

// CshFormatPath get csh format path declaration
func CshFormatPath() (output string) {
	return cshFormat("PATH")
}

// CshFormatManPath get csh format man path declaration
func CshFormatManPath() (output string) {
	return cshFormat("MANPATH")
}

func cshFormat(envVar string) (output string) {
	output = fmt.Sprintf("setenv %s=\"%s\";", envVar, strings.Join(finalPaths, ":"))

	return
}

// BashFormatPath get bash format path declaration
func BashFormatPath() (output string) {
	return bashFormat("PATH")
}

// BashFormatManPath get bash format man path declaration
func BashFormatManPath() (output string) {
	return bashFormat("MANPATH")
}

func bashFormat(envVar string) (output string) {
	output = fmt.Sprintf("%s=\"%s\"; export %s;", envVar, strings.Join(finalPaths, ":"), envVar)

	return
}

// extractPaths get list of paths from sources
func extractPaths(paths []string) (foundPaths []string) {
	for _, startingPath := range paths {
		home, err := os.UserHomeDir()
		if err != nil {
			os.Exit(1)
		}
		if strings.HasPrefix(startingPath, "~") {
			startingPath = strings.Replace(startingPath, "~", home, 1)
		}
		err = verifyPath(startingPath)
		if err != nil {
			continue
		}
		paths, err := getPathsInDir(startingPath)
		if err != nil {
			continue
		}
		for _, path := range paths {
			bytes, err := ioutil.ReadFile(path)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			scanner := bufio.NewScanner(strings.NewReader(string(bytes)))
			for scanner.Scan() {
				foundPaths = append(foundPaths, scanner.Text())
			}
			err = scanner.Err()
		}
	}

	return
}
