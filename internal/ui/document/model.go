package document

import (
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/13k/dm/internal/ui"
)

type Model struct {
	Styles Styles
	KeyMap KeyMap

	body    string
	bodyFmt string
	help    help.Model
}

func NewModel() Model {
	return Model{
		KeyMap: DefaultKeyMap(),
		Styles: DefaultStyles(),
		help:   help.NewModel(),
	}
}

func (m *Model) SetSize(w, _ int) {
	m.help.Width = w
}

func (m Model) Init() tea.Cmd { // nolint: gocritic
	log.Println("document.Init()")

	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) { // nolint: gocritic
	log.Printf("document.Update() -- [%T] %v", msg, msg)

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case *ui.DocumentRenderedMsg:
		m.body = msg.Body
		m.bodyFmt = msg.BodyColored
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.Save):
			cmd = m.save()
		case key.Matches(msg, m.KeyMap.Quit):
			cmd = tea.Quit
		}
	}

	return m, cmd
}

func (m *Model) save() tea.Cmd {
	return func() tea.Msg {
		return ui.NewSaveDocumentMsg(m.body)
	}
}

func (m Model) View() string { // nolint: gocritic
	var b strings.Builder

	b.WriteString(m.bodyView())
	b.WriteRune('\n')
	b.WriteString(m.helpView())

	return m.Styles.Frame.Render(b.String())
}

func (m *Model) bodyView() string {
	return m.Styles.BodyFrame.Render(m.bodyFmt)
}

func (m *Model) helpView() string {
	return m.Styles.HelpFrame.Render(m.help.View(m.KeyMap))
}
