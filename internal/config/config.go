package config

import (
	"io"

	"github.com/charmbracelet/glamour"
)

type Config struct {
	Output     io.Writer
	OutputPath string
	Renderer   *glamour.TermRenderer
}
