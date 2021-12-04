package app

import (
	"fmt"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/13k/dm/internal/config"
	"github.com/13k/dm/internal/ui"
	"github.com/13k/dm/internal/ui/document"
	"github.com/13k/dm/internal/ui/form"
)

var _ tea.Model = Model{}

type Model struct {
	Styles Styles

	cfg   *config.Config
	state appState
	form  form.Model
	doc   document.Model
	err   error
}

func NewModel(c *config.Config) Model {
	return Model{
		Styles: DefaultStyles(),
		cfg:    c,
		form:   form.NewModel(),
		doc:    document.NewModel(),
	}
}

func (m Model) Init() tea.Cmd { // nolint: gocritic
	log.Printf("app.Init() -- cfg: %#+v", m.cfg)

	return tea.Batch(
		tea.EnterAltScreen,
		m.form.Init(),
		m.loadInput(),
	)
}

func (m *Model) loadInput() tea.Cmd {
	if m.cfg.InputPath == "" {
		return nil
	}

	return ui.ParseDoc(m.cfg.InputPath)
}

func (m *Model) stateChange(s appState) {
	log.Printf("app.stateChange -- %s -> %s", m.state, s)
	m.state = s
}

func (m *Model) onResize(w, h int) tea.Cmd {
	log.Printf("app.onResize -- %d, %d", w, h)

	m.form.SetSize(w, h)
	m.doc.SetSize(w, h)

	return ui.NoopCmd
}

func (m *Model) onKeyPress(_ tea.KeyMsg) tea.Cmd {
	switch m.state {
	case stateError, stateDocumentSaved:
		return tea.Quit
	case stateShowForm, stateShowDocument:
	}

	return nil
}

func (m *Model) onError(err error) tea.Cmd {
	log.Printf("app.onError -- %v", err)

	m.err = err
	m.stateChange(stateError)

	return ui.NoopCmd
}

func (m *Model) onFormSubmitted(entries []string) tea.Cmd {
	log.Printf("app.onFormSubmitted -- %q", entries)

	return ui.RenderDoc(entries)
}

func (m *Model) onDocRendered(body, bodyColored string) tea.Cmd { // nolint: unparam
	log.Printf(
		"app.onDocRendered -- body size: %d, colored body size: %d",
		len(body),
		len(bodyColored),
	)

	m.stateChange(stateShowDocument)

	return nil
}

func (m *Model) onDocClipboard(body string) tea.Cmd {
	log.Printf("app.onDocClipboard -- body size: %d", len(body))

	return ui.ClipboardDoc(body)
}

func (m *Model) onDocSave(body string) tea.Cmd {
	log.Printf("app.onDocSave -- body size: %d", len(body))

	return ui.WriteDoc(body, m.cfg.OutputPath)
}

func (m *Model) onDocSaved() tea.Cmd {
	log.Printf("app.onDocSaved")

	m.stateChange(stateDocumentSaved)

	return ui.NoopCmd
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) { // nolint: gocritic
	if msg == nil {
		return m, nil
	}

	log.Printf("app.Update() -- [%T] %v", msg, msg)

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		cmd = m.onResize(msg.Width, msg.Height)
	case tea.KeyMsg:
		cmd = m.onKeyPress(msg)
	case *ui.ErrorMsg:
		cmd = m.onError(msg.Err)
	case *ui.FormSubmittedMsg:
		cmd = m.onFormSubmitted(msg.Entries)
	case *ui.DocumentRenderedMsg:
		cmd = m.onDocRendered(msg.Body, msg.BodyColored)
	case *ui.ClipboardDocumentMsg:
		cmd = m.onDocClipboard(msg.Body)
	case *ui.SaveDocumentMsg:
		cmd = m.onDocSave(msg.Body)
	case *ui.DocumentSavedMsg:
		cmd = m.onDocSaved()
	}

	// if one of the message handlers above returned non-nil command (including NoopCmd), return it.
	// otherwise delegate message handling to sub-models
	if cmd != nil {
		return m, cmd
	}

	return m.updateChildren(msg)
}

func (m *Model) updateChildren(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.state {
	case stateShowForm:
		m.form, cmd = m.form.Update(msg)

		return m, cmd
	case stateShowDocument:
		m.doc, cmd = m.doc.Update(msg)

		return m, cmd
	case stateError, stateDocumentSaved:
		return m, nil
	default:
		panic(fmt.Errorf("unknown app state %#+v", m.state))
	}
}

func (m Model) View() string { // nolint: gocritic
	switch m.state {
	case stateShowForm:
		m.form.Styles = m.Styles.Form

		return m.form.View()
	case stateShowDocument:
		m.doc.Styles = m.Styles.Document

		return m.doc.View()
	case stateError:
		return m.errorView()
	case stateDocumentSaved:
		return m.successView()
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

	msg := fmt.Sprintf("File saved to %s", m.cfg.OutputPath)

	b.WriteString(m.Styles.Success.Render(msg))
	b.WriteRune('\n')
	b.WriteString(m.Styles.Help.Render("Press any key to exit"))
	b.WriteRune('\n')

	return m.Styles.SuccessFrame.Render(b.String())
}
