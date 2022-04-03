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
var Trace *log.Logger
var once sync.Once
var verbose bool

func init() {
	once.Do(func() {
		Trace = log.New(io.Discard, "TRACE ", log.LUTC)
		Info = log.New(io.Discard, "INFO ", log.LUTC)
		Error = log.New(io.Discard, "ERROR ", log.LUTC)

		if args.Args.Verbose {
			Info = log.New(os.Stdout, "INFO ", log.LUTC)
			Error = log.New(os.Stdout, "ERROR ", log.LUTC)
		}
		if args.Args.Trace {
			Trace = log.New(os.Stdout, "TRACE ", log.LUTC)
		}
	})
}
