package paths

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/imarsman/pathhelper/cmd/args"
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

var homeDir, userPathsDirAbsolute, userManpathsDirAbsolute string

func init() {
	// Get user home dir
	homeDir, err := os.UserHomeDir()
	if err != nil {
		os.Exit(1)
	}

	// Define user paths and manpaths dirs
	userPathsDirAbsolute = filepath.Join(homeDir, ".config", "pathhelper", "paths.d")
	userManpathsDirAbsolute = filepath.Join(homeDir, ".config", "pathhelper", "manpaths.d")
}

// pathSet contains information on a path
type pathSet struct {
	kind       pathType
	systemPath string
	systemDir  string
	userDir    string
	paths      []string
	pathMap    map[string]struct{}
	mu         *sync.Mutex
}

func newPathSet(kind pathType, systemPath, systemDir, userDir string) (ps *pathSet) {
	ps = &pathSet{}
	ps.kind = kind
	ps.systemPath = systemPath
	ps.systemDir = systemDir
	ps.userDir = userDir
	ps.userDir = cleanDir(ps.userDir)
	ps.pathMap = make(map[string]struct{})
	ps.mu = new(sync.Mutex)

	return
}

// setup set up user directories if they don't exist
// Consider not using this from a flag as it's basically used every run
func setup() {
	fileMode := fs.FileMode(0755)

	if _, err := os.Stat(userPathsDirAbsolute); os.IsNotExist(err) {
		// path/to/whatever does not exist
		err = verifyPath(userPathsDirAbsolute)
		if err != nil {
			fmt.Fprintln(os.Stderr, "- creating user paths dir", userPathsDirAbsolute)
			err = os.MkdirAll(userPathsDirAbsolute, fileMode)
			if err != nil {
				fmt.Println(err)
			}
		}
	} else {
		_, err = checkUserOnlyRW(cleanDir(userPathDir))
		if err != nil {
			os.Chmod(userPathsDirAbsolute, 0700)
		}

	}
	if _, err := os.Stat(userManpathsDirAbsolute); os.IsNotExist(err) {
		err = verifyPath(userManpathsDirAbsolute)
		if err != nil {
			fmt.Fprintln(os.Stderr, "- creating user manpaths dir", userManpathsDirAbsolute)
			err = os.MkdirAll(userManpathsDirAbsolute, fileMode)
		}
	} else {
		_, err = checkUserOnlyRW(userManPathDir)
		if err != nil {
			os.Chmod(userManpathsDirAbsolute, 0700)
		}
	}

}

var configPaths *pathSet
var configManPaths *pathSet

func checkUserOnlyRW(path string) (userOnly bool, err error) {
	path = cleanDir(path)

	var info fs.FileInfo
	info, err = os.Stat(path)
	if err != nil {
		fmt.Printf("error getting file info: %v", err)
		return
	}
	mode := info.Mode()

	// Check for required permissions, i.e. rwx for user only
	//https://stackoverflow.com/questions/45429210/how-do-i-check-a-files-permissions-in-linux-using-go
	// https://codereview.stackexchange.com/questions/79020/bitwise-operators-for-permissions/79100#79100
	// Octal nontation
	if mode.Perm() != 0o0700 {
		err = fmt.Errorf("error for %s permissions must be 700 but are %v", path, mode)
		return
	}

	return
}

func init() {
	// Instantiate and populate - we can do this because the program runs once
	configPaths = newPathSet(pathTypePath, systemPathFile, systemPathDir, userPathDir)

	configManPaths = newPathSet(pathTypeManPath, systemManPathFile, systemManPathDir, userManPathDir)

	configPaths.populate()
	configManPaths.populate()
}

// verifyPath verify a path
func verifyPath(path string) (err error) {
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

			err = verifyPath(newPath)
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

// makeMapKey for hashing purposes remove all slashes all slashes, then lowercase
func makeMapKey(path string) string {
	path = strings.ReplaceAll(path, `/`, ``)
	path = strings.ToLower(path)

	return path
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
	logging.Trace.Printf("verify %s\n", file)

	lines := 0
	for scanner.Scan() {
		lines++

		path := strings.TrimSpace(scanner.Text())
		path = cleanDir(path)

		// who knows what might be encountered
		if lines > 99 {
			logging.Error.Printf("stopping processing of %s because line count max of 100", file)
			break
		}

		// skip comment lines
		if strings.HasPrefix(path, hash) {
			logging.Error.Printf("skipping line in file %s \"%s\"", filepath.Base(file), path)
			continue
		}
		// path has whitespace trimmed
		if len(path) == 0 {
			logging.Error.Printf("skipping empty path %s \"%s\"", filepath.Base(file), path)
			continue
		}
		// ~ is converted to $HOME by now so all paths should be absolute
		if path[0] != '/' {
			logging.Error.Printf("skipping path that does not begin with \"/\" %s \"%s\"", filepath.Base(file), path)
			continue
		}
		// Don't know what this would do but it seems never appropriate
		if path == "/" {
			logging.Error.Printf("skipping path that is root \"/\" %s \"%s\"", filepath.Base(file), path)
			continue
		}

		// By default skip verifying dirs
		// path_helper from Apple does not check dirs
		// Verifying adds about a ms to time to run
		// To find invalid paths run `pathhelper -z -V -v`
		if args.Args.Verify {
			logging.Info.Println("checking", path)

			// Check to ensure path is valid
			err = verifyPath(path)
			if err != nil {
				continue
			}
		}

		var sb strings.Builder
		var escaped = ""

		// Escape characters that can be bad for a shell to read in. This escaped value will be used for output. The
		// output of this program is an export statement for a shell variable.
		// Macos allows these characters in filenames and paths.
		// Similar logic to darwin C original
		// https://opensource.apple.com/source/shell_cmds/shell_cmds-162/path_helper/path_helper.c.auto.html
		for _, r := range path {
			switch r {
			case '"':
				logging.Error.Printf("escaping \" character in file %s \"%s\"", filepath.Base(file), path)
				sb.WriteRune('\\')
				sb.WriteRune(r)
				continue
			case '\'':
				logging.Error.Printf("escaping ' character in file %s \"%s\"", filepath.Base(file), path)
				sb.WriteRune('\\')
				sb.WriteRune(r)
				continue
			case '$':
				logging.Error.Printf("escaping $ character in file %s \"%s\"", filepath.Base(file), path)
				sb.WriteRune('\\')
				sb.WriteRune(r)
				continue
			case '\\':
				logging.Error.Printf("escaping \\ character in file %s \"%s\"", filepath.Base(file), path)
				sb.WriteRune('\\')
				sb.WriteRune(r)
				continue
			case '!':
				logging.Error.Printf("escaping ! character in file %s \"%s\"", filepath.Base(file), path)
				sb.WriteRune('\\')
				sb.WriteRune(r)
				continue
			}
			sb.WriteRune(r)
		}
		if sb.Len() > 0 {
			escaped = sb.String()
		}

		// Avoid duplicates by using map to keep track of what has been found so far
		// mutex used to protect any future concurrent access
		standardizedPath := makeMapKey(path)

		// Lock mutex for read and write block
		ps.mu.Lock()
		_, ok := ps.pathMap[standardizedPath]
		if ok {
			logging.Error.Printf("skipping duplicate path in file %s \"%s\"", filepath.Base(file), path)
			ps.mu.Unlock()
			continue
		}

		// It doesn't matter what we store. An empty struct consumes zero bytes.
		ps.pathMap[standardizedPath] = struct{}{}

		// Unlock mutex
		ps.mu.Unlock()

		// Use the escaped version for the PATH if it has passed tests
		if escaped != "" {
			path = escaped
		}

		// This is the only place we append to the paths list
		ps.paths = append(ps.paths, path)
	}
	logging.Trace.Printf("done %s in %v\n", file, time.Since(t1))

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

	// Run setup() to check user dirs before proceeding
	setup()

	var repeat = 80

	if args.Args.UserFirst {
		logging.Info.Println(strings.Repeat("-", repeat))
		// Get user paths.d entries
		t1 := time.Now()
		logging.Info.Printf("processing %s", ps.userDir)
		ps.addPathsFromDir(ps.userDir)
		logging.Info.Printf("processing %s took %v", ps.userDir, time.Since(t1).Milliseconds())
		logging.Info.Println(strings.Repeat("-", repeat))
	}

	if ps.kind != pathTypeManPath {
		logging.Info.Println(strings.Repeat("-", repeat))
	}
	t1 := time.Now()
	logging.Info.Printf("processing %s", ps.systemPath)
	// Get system path file lines
	ps.addPathsFromFile(ps.systemPath)
	logging.Info.Printf("processing %s took %v", ps.systemPath, time.Since(t1))

	logging.Info.Println(strings.Repeat("-", repeat))
	t1 = time.Now()
	logging.Info.Printf("processing %s", ps.systemDir)
	// Get system paths.d file entries
	ps.addPathsFromDir(ps.systemDir)
	logging.Info.Printf("processing %s took %v", ps.systemDir, time.Since(t1))

	logging.Info.Println(strings.Repeat("-", repeat))
	if !args.Args.UserFirst {
		t1 = time.Now()
		logging.Info.Printf("processing %s", ps.userDir)
		// Get user paths.d entries
		ps.addPathsFromDir(ps.userDir)
		logging.Info.Printf("processing %s took %v", ps.userDir, time.Since(t1))
		logging.Info.Println(strings.Repeat("-", repeat))
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

// Output paths in zsh format
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

// Output paths in csh format
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

// bashFormat get output in bash format
func (ps pathSet) bashFormat() (output string) {
	output = fmt.Sprintf("%s=\"%s\"; export %s;", ps.kind, strings.Join(ps.paths, ":"), ps.kind)

	return
}
