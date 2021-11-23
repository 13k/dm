package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/13k/dm/internal/app"
	"github.com/13k/dm/internal/config"
)

func init() {
	logPath := os.Getenv("DM_LOG")

	if logPath == "" {
		logPath = os.DevNull
	}

	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o666)
	if err != nil {
		panic(err)
	}

	log.SetOutput(f)
}

func main() {
	cfg := &config.Config{}

	if len(os.Args) > 1 {
		cfg.OutputFile = os.Args[1]
	}

	model := app.NewModel(cfg)

	if err := tea.NewProgram(model).Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Could not start program: %s\n", err)
		os.Exit(1)
	}
}
