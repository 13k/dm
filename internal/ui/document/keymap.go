package document

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

var _ help.KeyMap = &KeyMap{}

type KeyMap struct {
	Clipboard    key.Binding
	PublishSlack key.Binding
	Save         key.Binding
	Quit         key.Binding

	isSlackEnabled bool
}

func DefaultKeyMap() *KeyMap {
	return &KeyMap{
		Clipboard: key.NewBinding(
			key.WithKeys("y"),
			key.WithHelp("y", "copy to clipboard"),
		),
		PublishSlack: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "publish to slack channel"),
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
}

func (k *KeyMap) WithSlack(slack bool) *KeyMap {
	k.isSlackEnabled = slack

	return k
}

func (k *KeyMap) ShortHelp() []key.Binding {
	keys := []key.Binding{k.Clipboard}

	if k.isSlackEnabled {
		keys = append(keys, k.PublishSlack)
	}

	return append(keys, k.Save, k.Quit)
}

func (k *KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}
