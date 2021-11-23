package markdown

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/glamour"
)

const termRenderWidth = 80

var termRenderer *glamour.TermRenderer

func init() {
	var err error

	if termRenderer, err = newRenderer(termRenderWidth); err != nil {
		panic(fmt.Errorf("could not create markdown renderer: %w", err))
	}
}

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
	return termRenderer
}

func RenderTerm(body string) (string, error) {
	return TermRenderer().Render(body)
}
