package app

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/13k/dm/internal/config"
	"github.com/13k/dm/internal/message"
)

func saveDoc(cfg *config.Config, body string) tea.Cmd {
	log.Printf("saveDoc() -- config: %#+v", cfg)

	return func() tea.Msg {
		f, err := cfg.Output()
		if err != nil {
			log.Printf("saveDoc() -- open error: %v", err)
			return message.Error(err)
		}

		log.Printf("saveDoc() -- f: %q", f.Name())

		defer f.Close()

		n, err := f.WriteString(body)
		if err != nil {
			log.Printf("saveDoc() -- write error: %v", err)
			return message.Error(err)
		}

		log.Printf("saveDoc() -- written: %d", n)

		return message.DocumentSavedMsg{Filename: f.Name()}
	}
}
