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
		Foreground(Primary)

	Help = lipgloss.NewStyle().
		Foreground(Gray)

	Error = lipgloss.NewStyle().
		Background(Red).
		Padding(1, 1)

	EditorStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(Primary))

	CursorLineStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("57")).
			Foreground(lipgloss.Color("230"))
)
