package util

import (
	"os"
)

var (
	Cwd Path
)

func init() {
	cwd, err := os.Getwd()

	if err != nil {
		Fatal("could not determine current working directory: %w", err)
	}

	Cwd = NewPath(cwd)
}
