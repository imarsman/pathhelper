package logging

import (
	"fmt"
	"os"
)

func Error(parts ...any) {
	fmt.Fprintln(os.Stderr, parts...)
}
