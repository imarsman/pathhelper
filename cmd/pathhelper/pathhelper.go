package main

import (
	"github.com/alexflint/go-arg"
	"github.com/imarsman/pathhelper/cmd/paths"
)

// /usr/libexec/path_helper -s
// /etc/paths
// /etc/manpaths

var allPaths []string

func init() {

}

func main() {
	var args struct {
		Bash bool `arg:"-s,--bash"`
		CSH  bool `arg:"-c,--csh"`
		ZSH  bool `arg:"-z,--zsh"`
	}
	arg.MustParse(&args)

	paths.ExtractPaths()
}
