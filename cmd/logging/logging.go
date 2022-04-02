package logging

import (
	"io"
	"log"
	"os"
	"sync"

	"github.com/imarsman/pathhelper/cmd/args"
)

// Logger a simple logger for printing out info tied to reading files
var Logger *log.Logger
var once sync.Once
var verbose bool

func init() {
	once.Do(func() {
		if args.Args.Verbose {
			Logger = log.New(os.Stderr, "INFO ", log.LUTC)
		} else {
			Logger = log.New(io.Discard, "INFO ", log.LUTC)
		}
	})
}
