package form

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

var _ help.KeyMap = &KeyMap{}

type KeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Submit key.Binding
	Quit   key.Binding
}

func DefaultKeyMap() *KeyMap {
	return &KeyMap{
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
}

func (k *KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Submit, k.Quit}
}

func (k *KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}
