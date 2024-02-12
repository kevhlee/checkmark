package ui

import "github.com/charmbracelet/lipgloss"

var (
	styleCheck  = lipgloss.NewStyle().Foreground(lipgloss.Color("#afffd0"))
	styleCursor = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff80ee"))
	styleHeader = lipgloss.NewStyle().
			Bold(true).
			Underline(true).
			Foreground(lipgloss.Color("#acffee"))
)
