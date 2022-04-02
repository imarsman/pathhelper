package args

import "github.com/alexflint/go-arg"

func init() {
	arg.MustParse(&Args)
}

var Args struct {
	Bash    bool `arg:"-s,--bash"`
	CSH     bool `arg:"-c,--csh"`
	ZSH     bool `arg:"-z,--zsh"`
	Verbose bool `arg:"-v,--verbose"`
	Init    bool `arg:"-i,--init"`
}
