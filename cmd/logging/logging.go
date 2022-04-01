package logging

import (
	"fmt"
	"os"
)

// This logging implementation can change

// Error print error
func Error(parts ...any) {
	fmt.Fprintln(os.Stderr, parts...)
}
