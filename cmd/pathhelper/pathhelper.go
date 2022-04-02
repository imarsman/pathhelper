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

// The MacOS cli tool that inspired this
// /usr/libexec/path_helper

// setup set up user directories if they don't exist
func setup() {
	fmt.Println("setting up local directories")
	fmt.Println(strings.Repeat("-", len("setting up local directories")))
	fileMode := fs.FileMode(0755)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		os.Exit(1)
	}

	pathsDir := filepath.Join(homeDir, ".config", "pathhelper", "paths.d")
	manpathsDir := filepath.Join(homeDir, ".config", "pathhelper", "manpaths.d")

	err = paths.VerifyPath(pathsDir)
	if err != nil {
		fmt.Println("- creating user paths dir", pathsDir)
		err = os.MkdirAll(pathsDir, fileMode)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Printf("- user paths dir %s exists\n", pathsDir)
	}
	err = paths.VerifyPath(manpathsDir)
	if err != nil {
		fmt.Println("- creating user manpaths dir", manpathsDir)
		err = os.MkdirAll(manpathsDir, fileMode)
	} else {
		fmt.Printf("- user manpaths dir %s exists\n", manpathsDir)
	}
}

func main() {
	if args.Args.Init {
		setup()

		return
	}

	if args.Args.Bash {
		// Get bash style path and manpath setting
		fmt.Println(paths.BashFormatPath())
		fmt.Println(paths.BashFormatManPath())
	} else if args.Args.ZSH {
		// Get zsh style path and manpath setting
		fmt.Println(paths.ZshFormatPath())
		fmt.Println(paths.ZshFormatManPath())
	} else if args.Args.CSH {
		// Get csh style path and manpath setting
		fmt.Println(paths.CshFormatPath())
		fmt.Println(paths.CshFormatManPath())
	} else {
		// This is the default behaviour of /usr/libexec/path_helper
		fmt.Println(paths.BashFormatPath())
		fmt.Println(paths.BashFormatManPath())
	}
}
