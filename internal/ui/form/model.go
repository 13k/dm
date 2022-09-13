package form

import (
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/13k/dm/internal/ui"
	"github.com/13k/dm/internal/util"
)

const (
	defaultTitle = "Daily Meeting"
	buttonText   = "submit"
)

const (
	idxDone int = iota
	idxTodo
	idxBlocking
	idxSubmit
)

type inputConfig struct {
	prompt       string
	placeholder  string
	defaultValue string
}

var inputConfigs = make([]*inputConfig, ui.EntriesLen)

func init() {
	inputConfigs[idxDone] = &inputConfig{
		prompt:       "1. ",
		placeholder:  "done",
		defaultValue: "",
	}

	inputConfigs[idxTodo] = &inputConfig{
		prompt:       "2. ",
		placeholder:  "todo",
		defaultValue: "",
	}

	inputConfigs[idxBlocking] = &inputConfig{
		prompt:       "3. ",
		placeholder:  "blks",
		defaultValue: "none",
	}
}

type Model struct {
	Title  string
	KeyMap KeyMap
	Styles Styles

	focusIdx int
	inputs   []textinput.Model
	help     help.Model
}

func NewModel() Model {
	inputs := make([]textinput.Model, ui.EntriesLen)

	m := Model{
		Title:    defaultTitle,
		KeyMap:   DefaultKeyMap(),
		Styles:   DefaultStyles(),
		inputs:   inputs,
		focusIdx: idxDone,
		help:     help.NewModel(),
	}

	var t textinput.Model

	for i := range m.inputs {
		cfg := inputConfigs[i]
		t = textinput.NewModel()

		t.CursorStyle = m.Styles.Cursor
		t.Prompt = cfg.prompt
		t.Placeholder = cfg.placeholder

		if i == m.focusIdx {
			t.PromptStyle = m.Styles.FocusedPrompt
			t.TextStyle = m.Styles.FocusedText
		}

		inputs[i] = t
	}

	return m
}

func (m *Model) isSubmitFocused() bool {
	return m.focusIdx == idxSubmit
}

func (m *Model) SetSize(w, _ int) {
	m.help.Width = w
}

func (m Model) Init() tea.Cmd { //nolint: gocritic
	log.Println("form.Init()")

	return tea.Batch(
		textinput.Blink,
		m.inputs[m.focusIdx].Focus(),
	)
}

func (m *Model) onDocParsed(entries []string) {
	log.Printf("form.onDocParsed() -- entries: %q", entries)

	for i := range m.inputs {
		if len(entries) > i {
			m.inputs[i].SetValue(entries[i])
		}
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) { //nolint: gocritic
	log.Printf("form.Update() -- [%T] %v", msg, msg)

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case *ui.DocumentParsedMsg:
		m.onDocParsed(msg.Entries)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.KeyMap.Submit):
			if m.isSubmitFocused() {
				cmd = m.submit()
			} else {
				cmd = m.updateFocus(msg)
			}

			return m, cmd

		case key.Matches(msg, m.KeyMap.Up) || key.Matches(msg, m.KeyMap.Down):
			cmd = m.updateFocus(msg)

			return m, cmd
		}
	}

	cmd = m.updateInputs(msg)

	return m, cmd
}

func (m *Model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m *Model) updateFocus(msg tea.KeyMsg) tea.Cmd {
	if key.Matches(msg, m.KeyMap.Up) {
		m.focusIdx--
	} else {
		m.focusIdx++
	}

	if m.focusIdx > idxSubmit {
		m.focusIdx = idxDone
	} else if m.focusIdx < 0 {
		m.focusIdx = idxSubmit
	}

	var cmd tea.Cmd

	for i := range m.inputs {
		if i == m.focusIdx {
			cmd = m.inputs[i].Focus()

			m.inputs[i].PromptStyle = m.Styles.FocusedPrompt
			m.inputs[i].TextStyle = m.Styles.FocusedText
		} else {
			m.inputs[i].Blur()

			m.inputs[i].PromptStyle = m.Styles.BlurredPrompt
			m.inputs[i].TextStyle = m.Styles.BlurredText
		}
	}

	return cmd
}

func (m *Model) submit() tea.Cmd {
	return func() tea.Msg {
		entries := make([]string, len(m.inputs))

		for i := range m.inputs {
			val := m.inputs[i].Value()

			if val == "" {
				val = inputConfigs[i].defaultValue
			}

			entries[i] = val
		}

		return ui.NewFormSubmittedMsg(entries)
	}
}

func (m Model) View() string { //nolint: gocritic
	var b strings.Builder

	b.WriteString(m.formView())
	b.WriteRune('\n')
	b.WriteString(m.helpView())

	return b.String()
}

func (m *Model) formView() string {
	var b strings.Builder

	b.WriteString(m.headerView())
	b.WriteRune('\n')
	b.WriteString(m.inputsView())
	b.WriteRune('\n')
	b.WriteString(m.submitView())

	return m.Styles.Frame.Render(b.String())
}

func (m *Model) headerView() string {
	title := m.Styles.Title.Render(m.Title)
	date := m.Styles.Date.Render(util.TodayString())

	var b strings.Builder

	b.WriteString(title)
	b.WriteRune('\n')
	b.WriteString(date)

	return m.Styles.HeaderFrame.Render(b.String())
}

func (m *Model) inputsView() string {
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())

		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	return m.Styles.InputsFrame.Render(b.String())
}

func (m *Model) submitView() string {
	var b strings.Builder

	if m.isSubmitFocused() {
		b.WriteString(m.Styles.FocusedButton.Render(buttonText))
	} else {
		b.WriteString(m.Styles.BlurredButton.Render(buttonText))
	}

	return m.Styles.ButtonFrame.Render(b.String())
}

func (m *Model) helpView() string {
	var b strings.Builder

	b.WriteString(m.help.View(m.KeyMap))

	return m.Styles.HelpFrame.Render(b.String())
}
