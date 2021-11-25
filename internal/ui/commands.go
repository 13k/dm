package ui

import (
	"io"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/13k/dm/internal/markdown"
)

func NoopCmd() tea.Msg {
	return nil
}

func ParseDoc(path string) tea.Cmd {
	log.Printf("ParseDoc() -- create command")

	return func() tea.Msg {
		f, err := os.Open(path)
		if err != nil {
			log.Printf("ParseDoc() -- open file error: %v", err)
			return NewErrorMsg(err)
		}

		defer f.Close()

		src, err := io.ReadAll(f)
		if err != nil {
			log.Printf("ParseDoc() -- read error: %v", err)
			return NewErrorMsg(err)
		}

		var entries []string

		if items := markdown.ParseList(src); len(items) >= EntriesLen {
			entries = items[:EntriesLen]
		}

		return NewDocumentParsedMsg(entries)
	}
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

		return NewDocumentRenderedMsg(body, bodyColored)
	}
}

func WriteDoc(body, path string) tea.Cmd {
	log.Printf("WriteDoc() -- create command")

	return func() tea.Msg {
		f, err := os.Create(path)
		if err != nil {
			log.Printf("WriteDoc() -- create file error: %v", err)
			return NewErrorMsg(err)
		}

		defer f.Close()

		n, err := f.WriteString(body)
		if err != nil {
			log.Printf("WriteDoc() -- write error: %v", err)
			return NewErrorMsg(err)
		}

		log.Printf("WriteDoc() -- written: %d", n)

		return NewDocumentSavedMsg()
	}
}
