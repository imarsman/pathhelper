package main

import (
	"fmt"

	"github.com/imarsman/pathhelper/cmd/args"
	"github.com/imarsman/pathhelper/cmd/paths"
)

// /usr/libexec/path_helper -s
// /etc/paths
// /etc/manpaths

func init() {
}

var allPaths []string

// /usr/libexec/path_helper

func main() {
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
