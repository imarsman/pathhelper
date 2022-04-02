package logging

import (
	"io"
	"log"
	"os"
	"sync"

	"github.com/imarsman/pathhelper/cmd/args"
)

// Info a simple logger for printing out info tied to reading files
var Info *log.Logger
var Error *log.Logger
var once sync.Once
var verbose bool

func init() {
	once.Do(func() {
		if args.Args.Verbose {
			Info = log.New(os.Stderr, "INFO ", log.LUTC)
			Error = log.New(os.Stderr, "ERROR ", log.LUTC)
		} else {
			Info = log.New(io.Discard, "INFO ", log.LUTC)
			Error = log.New(io.Discard, "ERROR ", log.LUTC)
		}
	})
}
