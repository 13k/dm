package document

import (
	"github.com/charmbracelet/lipgloss"
)

var defaultStyles = Styles{
	Frame: lipgloss.NewStyle().
		Padding(2). //nolint: gomnd
		MarginBottom(1),
	MessageFrame: lipgloss.NewStyle().
		MarginTop(1),
	Message: lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")),
}

type Styles struct {
	// Widget styles
	Frame lipgloss.Style
	// Body styles
	BodyFrame lipgloss.Style
	// Message styles
	MessageFrame lipgloss.Style
	Message      lipgloss.Style
	// Help styles
	HelpFrame lipgloss.Style
}

func DefaultStyles() Styles {
	return defaultStyles
}
