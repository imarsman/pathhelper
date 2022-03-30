package paths

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var startingPaths = []string{"~/.config/pathhelper/paths.d",
	"/etc/paths",
	"/etc/paths.d"}
var startingManPaths = []string{"~/.config/pathhelper/manpaths.d",
	"/etc/manpaths",
	"/etc/manpaths.d"}
var finalPaths []string
var finalManPaths []string

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
				return
			}
			paths = append(paths, path)
		}
	}

	return
}

// ExtractPaths get list of paths from sources
func ExtractPaths() (paths []string, err error) {
	for _, path := range startingPaths {
		fmt.Println(path)
	}

	return
}

// ExtractManPaths get list of man paths from sources
func ExtractManPaths() (paths []string, err error) {

	for _, path := range startingPaths {
		fmt.Println(path)
	}

	return
}
