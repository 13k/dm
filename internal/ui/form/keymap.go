package form

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

var defaultKeyMap = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "shift+tab"),
		key.WithHelp("↑/<shift+tab>", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "tab"),
		key.WithHelp("↓/<tab>", "down"),
	),
	Submit: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("<enter>", "submit"),
	),
	Quit: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("<esc>", "quit"),
	),
}

var _ help.KeyMap = defaultKeyMap

type KeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Submit key.Binding
	Quit   key.Binding
}

func DefaultKeyMap() KeyMap {
	return defaultKeyMap
}

func (k KeyMap) ShortHelp() []key.Binding { // nolint: gocritic
	return []key.Binding{k.Up, k.Down, k.Submit, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding { // nolint: gocritic
	return [][]key.Binding{k.ShortHelp()}
}
