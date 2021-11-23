package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"

	"github.com/13k/dm/internal/config"
	"github.com/13k/dm/internal/ui/app"
	"github.com/13k/dm/internal/util"
)

const renderedDocWidth = 80

var (
	logPath string
	cwd     string
)

func init() {
	var err error

	cwd, err = os.Getwd()

	if err != nil {
		fatal("Could not determine current working directory: %v", err)
	}

	logPath = os.Getenv("DM_LOG")

	if logPath == "" {
		logPath = os.DevNull
	}

	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o666)
	if err != nil {
		fatal("Could not open log file %q: %v", logPath, err)
	}

	log.SetOutput(f)
}

func main() {
	var outputFile string

	if len(os.Args) > 1 {
		outputFile = os.Args[1]
	}

	output, err := createOutput(outputFile)
	if err != nil {
		fatal("Could not create output file: %v", err)
	}

	defer output.Close()

	renderer, err := newRenderer(renderedDocWidth)
	if err != nil {
		fatal("Could not create markdown renderer: %v", err)
	}

	cfg := &config.Config{
		Output:     output,
		OutputPath: output.Name(),
		Renderer:   renderer,
	}

	model := app.NewModel(cfg)

	if err := tea.NewProgram(model).Start(); err != nil {
		fatal("Could not start program: %v", err)
	}
}

func fatal(format string, args ...interface{}) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}

	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

func createOutput(filename string) (*os.File, error) {
	if filename == "" {
		basename := fmt.Sprintf("dm.%s.md", util.Today())
		filename = filepath.Join(cwd, basename)
	}

	return os.Create(filename)
}

func newRenderer(width int) (*glamour.TermRenderer, error) {
	return glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithEmoji(),
		glamour.WithWordWrap(width),
	)
}
