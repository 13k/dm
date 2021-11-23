package ui

import (
	"fmt"
	"io"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
)

func NoopCmd() tea.Msg {
	return nil
}

func RenderDoc(entries []string, renderer *glamour.TermRenderer) tea.Cmd {
	if len(entries) == 0 {
		log.Printf("RenderDoc() -- skip create command (empty entries)")
		return nil
	}

	log.Printf("RenderDoc() -- create command")

	return func() tea.Msg {
		log.Printf("RenderDoc() -- entries: %q", entries)

		var b strings.Builder

		for i, entry := range entries {
			fmt.Fprintf(&b, "%d. ", i+1)
			b.WriteString(entry)
			b.WriteRune('\n')
		}

		body := b.String()

		log.Printf("RenderDoc() -- body: %q", body)

		bodyColored, err := renderer.Render(body)
		if err != nil {
			log.Printf("RenderDoc() -- render error: %v", err)
			return NewErrorMsg(err)
		}

		log.Printf("RenderDoc() -- bodyColored size: %d", len(bodyColored))

		return DocumentRenderedMsg{
			Body:        body,
			BodyColored: bodyColored,
		}
	}
}

func SaveDoc(body string, output io.Writer) tea.Cmd {
	log.Printf("SaveDoc() -- create command")

	return func() tea.Msg {
		n, err := io.WriteString(output, body)
		if err != nil {
			log.Printf("SaveDoc() -- write error: %v", err)
			return NewErrorMsg(err)
		}

		log.Printf("SaveDoc() -- written: %d", n)

		return DocumentSavedMsg{}
	}
}
