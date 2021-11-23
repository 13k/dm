package config

import (
	"fmt"
	"os"

	"github.com/13k/dm/internal/util"
)

type Config struct {
	OutputFile string
}

// Output creates OutputFile if it's set, otherwise creates a temporary file and set OutputFile to
// the temporary file's path
func (c *Config) Output() (*os.File, error) {
	var (
		f   *os.File
		err error
	)

	if c.OutputFile != "" {
		f, err = os.Create(c.OutputFile)
	} else {
		pattern := fmt.Sprintf("dm.%s.*.md", util.Today())
		f, err = os.CreateTemp("", pattern)
	}

	return f, err
}
