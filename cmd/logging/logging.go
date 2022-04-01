package logging

import (
	"io/ioutil"
	"log"
	"os"
)

// SetVerbose set verbose output based on flag
func SetVerbose(v bool) {
	if !v {
		log.SetOutput(ioutil.Discard)
	} else {
		log.SetOutput(os.Stderr)
	}
}

// Log print error
func Log(parts ...any) {
	log.Println(parts...)
}
