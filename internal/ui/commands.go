package ui

import (
	"io"
	"log"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/13k/dm/internal/markdown"
)

func NoopCmd() tea.Msg {
	return nil
}

func RenderDoc(entries []string) tea.Cmd {
	if len(entries) == 0 {
		log.Printf("RenderDoc() -- skip create command (empty entries)")
		return nil
	}

	log.Printf("RenderDoc() -- create command")

	return func() tea.Msg {
		log.Printf("RenderDoc() -- entries: %q", entries)

		body := markdown.RenderList(entries)

		log.Printf("RenderDoc() -- body: %q", body)

		bodyColored, err := markdown.RenderTerm(body)
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

func WriteDoc(body string, output io.Writer) tea.Cmd {
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
