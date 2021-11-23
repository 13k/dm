package form

import (
	"github.com/charmbracelet/lipgloss"
)

var defaultStyles = Styles{
	Frame: lipgloss.NewStyle(),
	HeaderFrame: lipgloss.NewStyle().
		Margin(1, 1, 1, 3), // nolint: gomnd
	Title: lipgloss.NewStyle().
		Foreground(lipgloss.Color("63")),
	Date: lipgloss.NewStyle().
		Foreground(lipgloss.Color("5")),
	InputsFrame: lipgloss.NewStyle(),
	FocusedPrompt: lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("198")),
	FocusedText: lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")),
	BlurredPrompt: lipgloss.NewStyle(),
	BlurredText: lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{
			Light: "#909090",
			Dark:  "#828282",
		}),
	Cursor: lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")),
	ButtonFrame: lipgloss.NewStyle().
		Margin(1, 1, 1, 2), // nolint: gomnd
	FocusedButton: lipgloss.NewStyle().
		Bold(true).
		Padding(1).
		Foreground(lipgloss.Color("198")).
		Background(lipgloss.Color("234")),
	BlurredButton: lipgloss.NewStyle().
		Padding(1).
		Foreground(lipgloss.Color("205")),
}

type Styles struct {
	// Widget styles
	Frame lipgloss.Style
	// Header styles
	HeaderFrame lipgloss.Style
	Title       lipgloss.Style
	Date        lipgloss.Style
	// Inputs styles
	InputsFrame   lipgloss.Style
	FocusedPrompt lipgloss.Style
	FocusedText   lipgloss.Style
	BlurredPrompt lipgloss.Style
	BlurredText   lipgloss.Style
	Cursor        lipgloss.Style
	// Button styles
	ButtonFrame   lipgloss.Style
	FocusedButton lipgloss.Style
	BlurredButton lipgloss.Style
	// Help styles
	HelpFrame lipgloss.Style
}

func DefaultStyles() Styles {
	return defaultStyles
}
