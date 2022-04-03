package paths

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/imarsman/pathhelper/cmd/logging"
)

type pathType string

const (
	tilde                    = `~`
	hash                     = `#`
	pathTypePath    pathType = `PATH`
	pathTypeManPath pathType = `MANPATH`

	systemPathFile    = `/etc/paths`
	systemPathDir     = `/etc/paths.d`
	userPathDir       = `~/.config/pathhelper/paths.d`
	systemManPathFile = `/etc/manpaths`
	systemManPathDir  = `/etc/manpaths.d`
	userManPathDir    = `~/.config/pathhelper/manpaths.d`
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

var configPaths *pathSet
var configManPaths *pathSet

func init() {
	// Instantiate and populate - we can do this because the program runs once
	configPaths = newPathSet(pathTypePath, systemPathFile, systemPathDir, userPathDir)
	configManPaths = newPathSet(pathTypeManPath, systemManPathFile, systemManPathDir, userManPathDir)

	configPaths.populate()
	configManPaths.populate()
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

func (ps *pathSet) addPathsFromDir(path string) {
	// Get system file paths
	pathsInDir, err := filesInDir(path)
	if err != nil {
		logging.Info.Println(err)
		return
	}
	for i := 0; i < len(pathsInDir); i++ {
		ps.addPathsFromFile(pathsInDir[i])
	}

	return
}

// addPathsFromFile get valid paths from a file
func (ps *pathSet) addPathsFromFile(file string) {
	logging.Info.Println("checking", file)
	// The system path is a file with lines in it
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		logging.Info.Println(err)
		return
	}
	scanner := bufio.NewScanner(strings.NewReader(string(bytes)))
	t1 := time.Now()
	logging.Trace.Println("evaluating", file)
	for scanner.Scan() {
		path := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(path, hash) {
			logging.Error.Printf("skipping line in %s \"%s\"", filepath.Base(file), path)
			continue
		}
		path = cleanDir(path)
		logging.Info.Println("checking", path)
		err = VerifyPath(path)
		if err != nil {
			continue
		}
		// This is the only place we append to the paths list
		ps.paths = append(ps.paths, path)
	}
	logging.Trace.Printf("done %s in %v", file, time.Since(t1))
	err = scanner.Err()
	if err != nil {
		logging.Info.Println(err)
		return
	}

	return
}

// populate get paths in the order of system file, system dir, then user dir
// placing paths in front of system paths could be a security risk if the same
// named executable is in a part of the filesystem writeable by the user.
func (ps *pathSet) populate() (err error) {
	// The channel is used as a means of sending ordered data.
	// We intentionally do not want concurrency in channel add as we need to
	// maintain the ordering of the path variable we are building.

	// Get system path file lines
	logging.Info.Println("evaluating", ps.systemPath)
	ps.addPathsFromFile(ps.systemPath)

	// Get system paths.d file entries
	logging.Info.Println("evaluating", ps.systemDir)
	ps.addPathsFromDir(ps.systemDir)

	// Get user paths.d entries
	logging.Info.Println("evaluating", ps.userDir)
	ps.addPathsFromDir(ps.userDir)

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
