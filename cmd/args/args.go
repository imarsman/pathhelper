package args

import (
	"os"
	"strings"

	"github.com/alexflint/go-arg"
)

// This is put in a separate package to prevent circular dependencies between main and the logging package.

func init() {
	// Deal with go-args issue with testing
	testing := strings.HasSuffix(os.Args[0], ".test")
	if testing {
		p, _ := arg.NewParser(arg.Config{Program: "test"}, &Args)
		settings := []string{"--settings=test"}
		p.Parse(settings)
	} else {
		arg.MustParse(&Args)
	}
}

// Args args used in the app. Public for use in logging package.
var Args struct {
	Bash      bool `arg:"-s,--bash" help:"get bash format path settings"`
	ZSH       bool `arg:"-z,--zsh" help:"get zsh format path settings"`
	CSH       bool `arg:"-c,--csh" help:"get csh format path settings"`
	Verbose   bool `arg:"-v,--verbose" help:"display issues as paths evaluated"`
	Trace     bool `arg:"-t,--trace" help:"display very detailed activity"`
	Init      bool `arg:"-i,--init" help:"check and build user path dirs if necessary"`
	UserFirst bool `arg:"-u,--user-first" help:"put user directories first"`
}
