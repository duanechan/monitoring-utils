package model

import "github.com/charmbracelet/lipgloss"

// Colors
var (
	Primary   = lipgloss.Color("#6796BF")
	Secondary = lipgloss.Color("#FFFFFF")

	Red  = lipgloss.Color("#FF0000")
	Gray = lipgloss.Color("#808080")
)

var (
	Header = lipgloss.NewStyle().
		Foreground(Primary).
		Padding(3, 3)

	Error = lipgloss.NewStyle().
		Background(Red).
		Padding(1, 1)

	EditorStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(Primary))

	CursorLineStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("57")).
			Foreground(lipgloss.Color("230"))

	QuitButtonContainer = lipgloss.NewStyle().Width(14).Height(3).Align(lipgloss.Center, lipgloss.Center)

	QuitSelectedStyle = lipgloss.NewStyle().
				Padding(1, 2).
				Border(lipgloss.RoundedBorder()).
				Bold(true)
	QuitUnselectedStyle = lipgloss.NewStyle().
				Padding(2, 2).
				Foreground(Gray)
)
