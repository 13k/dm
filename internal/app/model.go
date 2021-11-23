package app

import (
	"fmt"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/13k/dm/internal/config"
	"github.com/13k/dm/internal/document"
	"github.com/13k/dm/internal/form"
	"github.com/13k/dm/internal/message"
)

const (
	stateShowForm state = iota
	stateShowDocument
)

type state int

func (s state) String() string {
	switch s {
	case stateShowForm:
		return "showForm"
	case stateShowDocument:
		return "showDocument"
	default:
		return "<unkown>"
	}
}

var _ tea.Model = Model{}

type Model struct {
	Styles Styles

	config    *config.Config
	state     state
	form      form.Model
	doc       document.Model
	err       error
	savedFile string
}

func NewModel(c *config.Config) Model {
	return Model{
		Styles: DefaultStyles(),
		config: c,
		form:   form.NewModel(),
		doc:    document.NewModel(),
	}
}

func (m *Model) stateChange(s state) {
	log.Printf("app.stateChange -- %s -> %s", m.state, s)
	m.state = s
}

func (m Model) Init() tea.Cmd { // nolint: gocritic
	return tea.Batch(
		tea.EnterAltScreen,
		m.form.Init(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) { // nolint: gocritic
	log.Printf("app.Update() -- [%T] %v", msg, msg)

	if _, ok := msg.(tea.KeyMsg); ok && (m.err != nil || m.savedFile != "") {
		return m, tea.Quit
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.form.SetSize(msg.Width, msg.Height)
		m.doc.SetSize(msg.Width, msg.Height)

		return m, nil
	case message.ErrorMsg:
		m.err = msg.Err

		return m, nil
	case message.FormSubmittedMsg:
		return m, renderDoc(msg.Entries)
	case message.DocumentRenderedMsg:
		m.stateChange(stateShowDocument)
	case message.SaveDocumentMsg:
		return m, saveDoc(m.config, msg.Body)
	case message.DocumentSavedMsg:
		m.savedFile = msg.Filename

		return m, nil
	}

	var cmd tea.Cmd

	switch m.state {
	case stateShowForm:
		m.form, cmd = m.form.Update(msg)

		return m, cmd
	case stateShowDocument:
		m.doc, cmd = m.doc.Update(msg)

		return m, cmd
	default:
		panic(fmt.Errorf("unknown app state %#+v", m.state))
	}
}

func (m Model) View() string { // nolint: gocritic
	if m.err != nil {
		return m.errorView()
	}

	if m.savedFile != "" {
		return m.successView()
	}

	switch m.state {
	case stateShowForm:
		m.form.Styles = m.Styles.Form

		return m.form.View()
	case stateShowDocument:
		m.doc.Styles = m.Styles.Document

		return m.doc.View()
	default:
		panic(fmt.Errorf("unknown app state %#+v", m.state))
	}
}

func (m *Model) errorView() string {
	var b strings.Builder

	msg := fmt.Sprintf("Error: %v", m.err)

	b.WriteString(m.Styles.Error.Render(msg))
	b.WriteRune('\n')
	b.WriteString(m.Styles.Help.Render("Press any key to exit"))
	b.WriteRune('\n')

	return m.Styles.ErrorFrame.Render(b.String())
}

func (m *Model) successView() string {
	var b strings.Builder

	msg := fmt.Sprintf("File saved to %s", m.savedFile)

	b.WriteString(m.Styles.Success.Render(msg))
	b.WriteRune('\n')
	b.WriteString(m.Styles.Help.Render("Press any key to exit"))
	b.WriteRune('\n')

	return m.Styles.SuccessFrame.Render(b.String())
}
