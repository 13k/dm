package util

import (
	"os"
)

var (
	Cwd string
)

func init() {
	var err error

	Cwd, err = os.Getwd()

	if err != nil {
		Fatal("could not determine current working directory: %w", err)
	}
}
