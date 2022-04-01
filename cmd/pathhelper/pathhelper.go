package main

import (
	"fmt"

	"github.com/alexflint/go-arg"
	"github.com/imarsman/pathhelper/cmd/paths"
)

// /usr/libexec/path_helper -s
// /etc/paths
// /etc/manpaths

var allPaths []string

// /usr/libexec/path_helper

var args struct {
	Bash bool `arg:"-s,--bash"`
	CSH  bool `arg:"-c,--csh"`
	ZSH  bool `arg:"-z,--zsh"`
}

func main() {
	arg.MustParse(&args)

	if args.Bash {
		fmt.Println(paths.BashFormatPath())
		fmt.Println(paths.BashFormatManPath())
	} else if args.ZSH {
		fmt.Println(paths.ZshFormatPath())
		fmt.Println(paths.ZshFormatManPath())
	} else if args.CSH {
		fmt.Println(paths.CshFormatPath())
		fmt.Println(paths.CshFormatManPath())
	} else {
		fmt.Println(paths.BashFormatPath())
		fmt.Println(paths.BashFormatManPath())
	}
}
