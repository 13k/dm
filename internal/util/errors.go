package util

import (
	"fmt"
	"os"
	"strings"
)

func Fatalf(format string, args ...interface{}) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}

	fmt.Fprintf(os.Stderr, format, args...)

	os.Exit(1)
}

func Must(err error) {
	if err != nil {
		Fatalf("Error: %v", err)
	}
}
