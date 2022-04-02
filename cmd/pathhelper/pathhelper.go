package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/imarsman/pathhelper/cmd/args"
	"github.com/imarsman/pathhelper/cmd/paths"
)

// The cli tool that inspired this
// /usr/libexec/path_helper

// func verifyPath(path string) (err error) {
// 	if _, err = os.Stat(path); err != nil {
// 		logging.Logger.Println(err)
// 		return
// 	}

// 	return
// }

func setup() {
	fmt.Println("setting up local directories")
	fmt.Println(strings.Repeat("-", len("setting up local directories")))
	dirMode := fs.FileMode(int(0777))

	homeDir, err := os.UserHomeDir()
	if err != nil {
		os.Exit(1)
	}
	configDirPath := filepath.Join(homeDir, ".config")
	err = paths.VerifyPath(configDirPath)
	if err != nil {
		fmt.Println("Creating ~/.config")
		err = os.MkdirAll(configDirPath, dirMode)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		fmt.Println("found", configDirPath)
	}
	pathhelperDirPath := filepath.Join(configDirPath, "pathhelper")
	err = paths.VerifyPath(pathhelperDirPath)
	if err != nil {
		fmt.Println("Creating ~/.config/pathhelper")
		err = os.MkdirAll(pathhelperDirPath, dirMode)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		fmt.Println("found", pathhelperDirPath)
	}
	pathsDirPath := filepath.Join(pathhelperDirPath, "paths.d")
	err = paths.VerifyPath(pathsDirPath)
	if err != nil {
		fmt.Println("Creating ~/.config/pathhelper/pathsPath.d")
		err = os.MkdirAll(pathsDirPath, dirMode)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println("found", pathsDirPath)
	}
	manpathsDirPath := filepath.Join(pathhelperDirPath, "manpaths.d")
	err = paths.VerifyPath(manpathsDirPath)
	if err != nil {
		fmt.Println("Creating ~/.config/pathhelper/manpaths.d")
		err = os.MkdirAll(manpathsDirPath, dirMode)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		fmt.Println("found", manpathsDirPath)
	}
}

var allPaths []string

func main() {
	if args.Args.Init {
		setup()

		return
	}

	if args.Args.Bash {
		fmt.Println(paths.BashFormatPath())
		fmt.Println(paths.BashFormatManPath())
	} else if args.Args.ZSH {
		fmt.Println(paths.ZshFormatPath())
		fmt.Println(paths.ZshFormatManPath())
	} else if args.Args.CSH {
		fmt.Println(paths.CshFormatPath())
		fmt.Println(paths.CshFormatManPath())
	} else {
		fmt.Println(paths.BashFormatPath())
		fmt.Println(paths.BashFormatManPath())
	}
}
