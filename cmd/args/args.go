package args

import "github.com/alexflint/go-arg"

func init() {
	arg.MustParse(&Args)
}

// Args args used in the app. Public for use in logging package.
var Args struct {
	Bash    bool `arg:"-s,--bash" help:"get bash format path settings"`
	ZSH     bool `arg:"-z,--zsh" help:"get zsh format path settings"`
	CSH     bool `arg:"-c,--csh" help:"get csh format path settings"`
	Verbose bool `arg:"-v,--verbose" help:"display issues as paths evaluated"`
	Init    bool `arg:"-i,--init" help:"check and build user path dirs if necessary"`
}
