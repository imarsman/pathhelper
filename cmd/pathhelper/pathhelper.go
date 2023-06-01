package main

import (
	"fmt"

	"github.com/imarsman/pathhelper/cmd/args"
	"github.com/imarsman/pathhelper/cmd/paths"
)

// The MacOS cli tool that inspired this
// /usr/libexec/path_helper

func main() {
	// if args.Args.Init {
	// 	paths.Setup()

	// 	return
	// }

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
