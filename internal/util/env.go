package util

import (
	"os"
)

func Cwd() Path {
	cwd, err := os.Getwd()
	if err != nil {
		Fatalf("could not determine current working directory: %w", err)
	}

	return NewPath(cwd)
}
