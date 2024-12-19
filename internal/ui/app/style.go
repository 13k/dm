package app

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/13k/dm/internal/ui/document"
	"github.com/13k/dm/internal/ui/form"
)

type Styles struct {
	Form         form.Styles
	Document     document.Styles
	SuccessFrame lipgloss.Style
	Success      lipgloss.Style
	ErrorFrame   lipgloss.Style
	Error        lipgloss.Style
	Help         lipgloss.Style
}

func DefaultStyles() Styles {
	styles := Styles{
		SuccessFrame: lipgloss.NewStyle().
			Margin(1).
			Padding(1),
		Success: lipgloss.NewStyle().
			Margin(1, 0).
			Padding(1).
			Foreground(lipgloss.Color("205")).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("205")),
		ErrorFrame: lipgloss.NewStyle().
			Margin(1).
			Padding(1),
		Error: lipgloss.NewStyle().
			Margin(1, 0).
			Padding(1).
			Foreground(lipgloss.Color("197")).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("197")),
		Help: lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{
				Light: "#909090",
				Dark:  "#626262",
			}),
	}

	docStyles := document.DefaultStyles()
	formStyles := form.DefaultStyles()

	helpFrameStyle := lipgloss.NewStyle().
		Margin(1, 1, 0, 1).
		PaddingTop(1).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.AdaptiveColor{
			Light: "#909090",
			Dark:  "#626262",
		}).
		BorderTop(true)

	styles.Document = docStyles
	styles.Document.HelpFrame = helpFrameStyle
	styles.Form = formStyles
	styles.Form.HelpFrame = helpFrameStyle

	return styles
}
