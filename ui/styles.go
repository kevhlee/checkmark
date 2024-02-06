package ui

import "github.com/charmbracelet/lipgloss"

var (
	styleCheck  = lipgloss.NewStyle().Foreground(lipgloss.Color("#afffd0"))
	styleCursor = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff80ee"))

	styleHeader = lipgloss.NewStyle().
			Bold(true).
			Underline(true).
			Foreground(lipgloss.Color("#acffee"))

	styleConfirmSelected = lipgloss.NewStyle().
				Margin(1, 1).
				Padding(0, 3).
				Background(lipgloss.Color("212")).
				Foreground(lipgloss.Color("230"))

	styleConfirmUnselected = lipgloss.NewStyle().
				Margin(1, 1).
				Padding(0, 3).
				Background(lipgloss.Color("235")).
				Foreground(lipgloss.Color("254"))

	styleEditorSelected   = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff"))
	styleEditorUnselected = lipgloss.NewStyle().Foreground(lipgloss.Color("#7f7f7f"))
)
