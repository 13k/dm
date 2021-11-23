package app

import (
	"fmt"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"

	"github.com/13k/dm/internal/message"
)

const coloredBodyWidth = 80

var renderer *glamour.TermRenderer

func init() {
	var err error

	renderer, err = glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithEmoji(),
		glamour.WithWordWrap(coloredBodyWidth),
	)

	if err != nil {
		panic(err)
	}
}

func renderDoc(entries []string) tea.Cmd {
	log.Printf("renderDoc() -- entries: %q", entries)

	if len(entries) == 0 {
		return nil
	}

	return func() tea.Msg {
		var b strings.Builder

		for i, entry := range entries {
			fmt.Fprintf(&b, "%d. ", i+1)
			b.WriteString(entry)
			b.WriteRune('\n')
		}

		body := b.String()

		log.Printf("renderDoc() -- body: %q", body)

		bodyColored, err := renderer.Render(body)
		if err != nil {
			log.Printf("renderDoc() -- render error: %v", err)
			return message.Error(err)
		}

		log.Println("renderDoc() -- bodyColored")

		return message.DocumentRenderedMsg{
			Body:        body,
			BodyColored: bodyColored,
		}
	}
}
