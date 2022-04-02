package paths

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/imarsman/pathhelper/cmd/logging"
)

var configPaths *pathSet
var configManPaths *pathSet

func init() {
	configPaths = newPathSet(pathPath, "/etc/paths", "/etc/paths.d", "~/.config/pathhelper/paths.d")
	configManPaths = newPathSet(manPath, "/etc/manpaths", "/etc/manpaths.d", "~/.config/pathhelper/manpaths.d")

	configPaths.populate()
	configManPaths.populate()
}

type pathType string

const (
	tilde             = `~`
	hash              = `#`
	pathPath pathType = `PATH`
	manPath  pathType = `MANPATH`
)

type pathSet struct {
	kind       pathType
	systemPath string
	systemDir  string
	userDir    string
	paths      []string
}

func newPathSet(kind pathType, systemPath, systemDir, userDir string) (ps *pathSet) {
	ps = &pathSet{}
	ps.kind = kind
	ps.systemPath = systemPath
	ps.systemDir = systemDir
	ps.userDir = userDir
	ps.userDir = cleanDir(ps.userDir)

	return
}

// VerifyPath verify a path
func VerifyPath(path string) (err error) {
	// Check if file exists and return with error if not
	if _, err = os.Stat(path); errors.Is(err, os.ErrNotExist) {
		logging.Error.Println(err)
		return
	}
	// Check that the file can be opened as read only
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		if os.IsPermission(err) {
			logging.Info.Printf("Unable to read %s\n", path)
			return
		}

	}

	return
}

var wg sync.WaitGroup

// populate get paths in the order of system file, system dir, then user dir
// placing paths in front of system paths could be a security risk if the same
// named executable is in a part of the filesystem writeable by the user.
func (ps *pathSet) populate() (err error) {
	var c = make(chan (string))

	go func() {
		// wg.Add(1)
		for {
			path, ok := <-c
			if !ok {
				// fmt.Println(ok)
				break
			}
			ps.paths = append(ps.paths, path)
		}
	}()

	logging.Info.Println("evaluating", ps.systemPath)

	// Get system path file lines
	addPathsFromFile(ps.systemPath, c)

	logging.Info.Println("evaluating", ps.systemDir)
	addPathsFromDir(ps.systemDir, c)

	logging.Info.Println("evaluating", ps.userDir)
	addPathsFromDir(ps.userDir, c)

	wg.Wait()
	close(c)

	// for line := range c {
	// 	ps.paths = append(ps.paths, line)
	// }

	return
}

func addPathsFromDir(path string, c chan (string)) (lines []string) {
	// Get system file paths
	pathsInDir, err := filesInDir(path)
	if err != nil {
		logging.Info.Println(err)
		return
	}
	for i := 0; i < len(pathsInDir); i++ {
		addPathsFromFile(pathsInDir[i], c)
	}

	return
}

// addPathsFromFile get valid paths from a file
func addPathsFromFile(file string, c chan (string)) {
	wg.Add(1)
	defer wg.Done()
	// The system path is a file with lines in it
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		logging.Info.Println(err)
		return
	}
	scanner := bufio.NewScanner(strings.NewReader(string(bytes)))
	for scanner.Scan() {
		path := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(path, hash) {
			continue
		}
		path = cleanDir(path)
		logging.Info.Println("checking", path)
		err = VerifyPath(path)
		if err != nil {
			continue
		}
		c <- path
	}
	err = scanner.Err()
	if err != nil {
		logging.Info.Println(err)
		return
	}

	return
}

func cleanDir(path string) (cleanDir string) {
	var homeDir, err = os.UserHomeDir()
	if err != nil {
		logging.Info.Println(err)
		return
	}

	cleanDir = strings.TrimSpace(path)
	if strings.HasPrefix(path, tilde) {
		cleanDir = cleanDir[1:]
		cleanDir = filepath.Join(homeDir, cleanDir)
	}

	return
}

// filesInDir get list of valid files in a dir
func filesInDir(basePath string) (paths []string, err error) {
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		logging.Info.Println(err)
		logging.Info.Println(err)
		return
	}

	for _, file := range files {
		if !file.IsDir() {
			var newPath = filepath.Join(basePath, file.Name())
			err = VerifyPath(newPath)
			if err != nil {
				logging.Info.Printf("can't read %s %v", newPath, err)
				continue
			}
			paths = append(paths, newPath)
		}
	}

	return
}

// Paths get list of paths
func Paths() (paths []string) {
	return configPaths.paths
}

// ManPaths get list of man paths
func ManPaths() (paths []string) {
	return configManPaths.paths
}

// ZshFormatPath get zsh format path declaration
func ZshFormatPath() (output string) {
	return configPaths.zshFormat()
}

// ZshFormatManPath get zsh format man path declaration
func ZshFormatManPath() (output string) {
	return configManPaths.zshFormat()
}

func (ps pathSet) zshFormat() (output string) {
	output = fmt.Sprintf("export %s=\"%s\";", ps.kind, strings.Join(ps.paths, ":"))
	return
}

// CshFormatPath get csh format path declaration
func CshFormatPath() (output string) {
	return configPaths.cshFormat()
}

// CshFormatManPath get csh format man path declaration
func CshFormatManPath() (output string) {
	return configManPaths.cshFormat()
}

func (ps pathSet) cshFormat() (output string) {
	output = fmt.Sprintf("setenv %s=\"%s\";", ps.kind, strings.Join(ps.paths, ":"))

	return
}

// BashFormatPath get bash format path declaration
func BashFormatPath() (output string) {
	return configPaths.bashFormat()
}

// BashFormatManPath get bash format man path declaration
func BashFormatManPath() (output string) {
	return configManPaths.bashFormat()
}

func (ps pathSet) bashFormat() (output string) {
	output = fmt.Sprintf("%s=\"%s\"; export %s;", ps.kind, strings.Join(ps.paths, ":"), ps.kind)

	return
}
