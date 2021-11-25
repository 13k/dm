package util

import (
	"fmt"
	"os"
	"strings"
)

func Fatal(format string, args ...interface{}) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}

	fmt.Fprintf(os.Stderr, format, args...)

	os.Exit(1)
}
