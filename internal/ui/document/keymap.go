package document

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

var defaultKeyMap = KeyMap{
	Clipboard: key.NewBinding(
		key.WithKeys("y"),
		key.WithHelp("y", "copy to clipboard"),
	),
	Save: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "save and quit"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc"),
		key.WithHelp("q/<esc>", "quit without saving"),
	),
}

var _ help.KeyMap = defaultKeyMap

type KeyMap struct {
	Clipboard key.Binding
	Save      key.Binding
	Quit      key.Binding
}

func DefaultKeyMap() KeyMap {
	return defaultKeyMap
}

func (k KeyMap) ShortHelp() []key.Binding { // nolint: gocritic
	return []key.Binding{k.Clipboard, k.Save, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding { // nolint: gocritic
	return [][]key.Binding{k.ShortHelp()}
}
