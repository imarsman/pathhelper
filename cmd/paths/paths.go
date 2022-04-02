package paths

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

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
	if _, err = os.Stat(path); err != nil {
		logging.Logger.Println(err)
		return
	}

	return
}

func (ps *pathSet) populate() (err error) {
	// Get user dir paths
	pathsInDir, err := filesInDir(ps.userDir)
	if err != nil {
		logging.Logger.Println(err)
		return
	}
	for _, p := range pathsInDir {
		lines := pathsFromFile(p)
		for i, line := range lines {
			lines[i] = cleanDir(line)
		}
		ps.paths = append(ps.paths, lines...)
	}

	// Get system path file lines
	lines := pathsFromFile(ps.systemPath)
	for i, line := range lines {
		lines[i] = cleanDir(line)
	}
	ps.paths = append(ps.paths, lines...)

	// Get system file paths
	pathsInDir, err = filesInDir(ps.systemDir)
	if err != nil {
		logging.Logger.Println(err)
		return
	}
	for _, p := range pathsInDir {
		lines := pathsFromFile(p)
		for i, line := range lines {
			lines[i] = cleanDir(line)
		}
		ps.paths = append(ps.paths, lines...)
	}

	return
}

// pathsFromFile get valid paths from a file
func pathsFromFile(file string) (lines []string) {
	// The system path is a file with lines in it
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		logging.Logger.Println(err)
		return
	}
	scanner := bufio.NewScanner(strings.NewReader(string(bytes)))
	for scanner.Scan() {
		path := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(path, hash) {
			continue
		}
		path = cleanDir(path)
		err = VerifyPath(path)
		if err != nil {
			continue
		}
		lines = append(lines, scanner.Text())
	}
	err = scanner.Err()
	if err != nil {
		logging.Logger.Println(err)
		return
	}

	return
}

func cleanDir(path string) (cleanDir string) {
	var homeDir, err = os.UserHomeDir()
	if err != nil {
		logging.Logger.Println(err)
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
		logging.Logger.Println(err)
		return
	}

	for _, file := range files {
		if !file.IsDir() {
			newPath := filepath.Join(basePath, file.Name())
			err = VerifyPath(newPath)
			if err != nil {
				logging.Logger.Println(err)
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
