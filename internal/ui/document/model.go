package document

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/13k/dm/internal/config"
	"github.com/13k/dm/internal/ui"
)

const (
	iconClipboard = "📋"
	iconSlack     = "📨"
)

type Model struct {
	Styles Styles
	KeyMap KeyMap

	cfg  *config.Config
	help help.Model

	body             string
	bodyFmt          string
	clipboardWritten *ui.ClipboardWrittenMsg
	slackSent        *ui.PublishedSlackMsg
}

func NewModel(cfg *config.Config) Model {
	return Model{
		KeyMap: DefaultKeyMap().WithSlack(cfg.SlackChannel != ""),
		Styles: DefaultStyles(),
		cfg:    cfg,
		help:   help.NewModel(),
	}
}

func (m *Model) SetSize(w, _ int) {
	m.help.Width = w
}

func (m Model) Init() tea.Cmd { //nolint: gocritic
	log.Println("document.Init()")

	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) { //nolint: gocritic
	log.Printf("document.Update() -- [%T] %v", msg, msg)

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case *ui.DocumentRenderedMsg:
		m.body = msg.Body
		m.bodyFmt = msg.BodyColored
	case *ui.ClipboardWrittenMsg:
		m.clipboardWritten = msg
	case *ui.PublishedSlackMsg:
		m.slackSent = msg
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.Clipboard):
			cmd = m.copyToClipboard()
		case key.Matches(msg, m.KeyMap.PublishSlack):
			cmd = m.publishSlack()
		case key.Matches(msg, m.KeyMap.Save):
			cmd = m.save()
		case key.Matches(msg, m.KeyMap.Quit):
			cmd = tea.Quit
		}
	}

	return m, cmd
}

func (m *Model) copyToClipboard() tea.Cmd {
	return func() tea.Msg {
		return ui.NewClipboardDocumentMsg(m.body)
	}
}

func (m *Model) publishSlack() tea.Cmd {
	return func() tea.Msg {
		return ui.NewPublishSlackMsg(m.cfg.SlackChannel, m.body)
	}
}

func (m *Model) save() tea.Cmd {
	return func() tea.Msg {
		return ui.NewSaveDocumentMsg(m.body)
	}
}

func (m Model) View() string { //nolint: gocritic
	var b strings.Builder

	b.WriteString(m.bodyView())
	b.WriteRune('\n')
	b.WriteString(m.messageView())
	b.WriteRune('\n')
	b.WriteString(m.helpView())

	return m.Styles.Frame.Render(b.String())
}

func (m *Model) bodyView() string {
	return m.Styles.BodyFrame.Render(m.bodyFmt)
}

func (m *Model) messageView() string {
	var b strings.Builder

	if m.clipboardWritten != nil {
		msg := fmt.Sprintf("%s  copied to the clipboard", iconClipboard)

		b.WriteString(m.Styles.Message.Render(msg))
		b.WriteRune('\n')
	}

	if m.slackSent != nil {
		msg := fmt.Sprintf(
			"%s  sent to slack channel %s (%s)",
			iconSlack,
			m.slackSent.Channel,
			m.slackSent.Timestamp,
		)

		b.WriteString(m.Styles.Message.Render(msg))
		b.WriteRune('\n')
	}

	var view string

	body := b.String()

	if body != "" {
		view = m.Styles.MessageFrame.Render(body)
	}

	return view
}

func (m *Model) helpView() string {
	return m.Styles.HelpFrame.Render(m.help.View(m.KeyMap))
}
