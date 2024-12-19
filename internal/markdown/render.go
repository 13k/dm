package markdown

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/glamour"
)

const (
	termRenderWidth = 80
)

func newRenderer(width int) (*glamour.TermRenderer, error) {
	return glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithEmoji(),
		glamour.WithWordWrap(width),
	)
}

func RenderList(entries []string) string {
	if len(entries) == 0 {
		return ""
	}

	var b strings.Builder

	for i, entry := range entries {
		fmt.Fprintf(&b, "%d. ", i+1)
		b.WriteString(entry)
		b.WriteRune('\n')
	}

	return b.String()
}

func TermRenderer() *glamour.TermRenderer {
	termRenderer, err := newRenderer(termRenderWidth)
	if err != nil {
		panic(fmt.Errorf("could not create markdown renderer: %w", err))
	}

	return termRenderer
}

func RenderTerm(body string) (string, error) {
	return TermRenderer().Render(body)
}
