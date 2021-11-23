package config

import (
	"io"
)

type Config struct {
	Output     io.Writer
	OutputPath string
}
