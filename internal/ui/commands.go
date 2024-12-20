package ui

import (
	"io"
	"log"
	"os"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/slack-go/slack"

	"github.com/13k/dm/internal/keyring"
	"github.com/13k/dm/internal/markdown"
	"github.com/13k/dm/internal/util"
)

const (
	slackAuthTokenKey = "slack_auth_token" //nolint:gosec // it's a key name, not a secret key
)

func NoopCmd() tea.Msg { //nolint:ireturn // bubbletea command func
	return nil
}

func ParseDoc(path util.Path) tea.Cmd {
	log.Printf("ParseDoc() -- create command")

	return func() tea.Msg {
		f, err := os.Open(path.String())
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

func ClipboardDoc(body string) tea.Cmd {
	log.Printf("ClipboardDoc() -- create command")

	return func() tea.Msg {
		err := clipboard.WriteAll(body)
		if err != nil {
			log.Printf("ClipboardDoc() -- clipboard write error: %v", err)

			return NewErrorMsg(err)
		}

		return NewClipboardWrittenMsg()
	}
}

func PublishSlackDoc(channel, body string) tea.Cmd {
	log.Printf("PublishSlackDoc() -- create command")

	return func() tea.Msg {
		authToken, err := keyring.Get(slackAuthTokenKey)
		if err != nil {
			log.Printf("PublishSlackDoc() -- keyring error: %v", err)

			return NewErrorMsg(err)
		}

		client := slack.New(authToken)
		msgOpts := []slack.MsgOption{
			slack.MsgOptionAsUser(true),
			slack.MsgOptionText(body, false),
		}

		resChannel, resTimestamp, err := client.PostMessage(channel, msgOpts...)
		if err != nil {
			log.Printf("PublishSlackDoc() -- slack error: %v", err)

			return NewErrorMsg(err)
		}

		return NewPublishedSlackMsg(resChannel, resTimestamp)
	}
}

func WriteDoc(body string, path util.Path) tea.Cmd {
	log.Printf("WriteDoc() -- create command")

	return func() tea.Msg {
		f, err := os.Create(path.String())
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
