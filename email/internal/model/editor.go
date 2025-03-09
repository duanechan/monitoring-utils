package model

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	email "github.com/duanechan/monitoring-utils/email/internal"
)

type (
	EditorModel struct {
		enabled bool
		records table.Model
	}
)

func initEditor(result email.ParseResult) table.Model {
	rows := []table.Row{}

	redStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	greenStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("46"))

	for i, r := range result.Raw {
		var status string
		if e, exists := result.BadEmails[i+1]; exists {
			status = redStyle.Render(fmt.Sprintf("✖ %s", e))
		} else {
			status = greenStyle.Render("✔")
		}

		rows = append(rows, table.Row{fmt.Sprintf("%d", i+1), r[0], r[1], status})
	}

	cols := []table.Column{
		{Title: "Row", Width: 5},
		{Title: "Name", Width: 25},
		{Title: "Email", Width: 25},
		{Title: "Valid", Width: 100}, // Apply color to header
	}

	t := table.New(
		table.WithHeight(len(result.Raw)+1),
		table.WithColumns(cols),
		table.WithRows(rows),
		table.WithStyles(table.Styles{
			Header:   lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("210")), // Default headers
			Selected: lipgloss.NewStyle().Background(lipgloss.Color("236")).Foreground(lipgloss.Color("15")),
		}),
	)

	return t
}

func (e EditorModel) Init() tea.Cmd {
	return nil
}

func (e EditorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			e.records.MoveUp(0)
		case "down":
			e.records.MoveDown(0)
		}
	}

	if e.enabled {
		e.records, cmd = e.records.Update(msg)
	}
	return e, cmd
}

func (e EditorModel) View() string {
	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63")).
		Render(e.records.View())
}
