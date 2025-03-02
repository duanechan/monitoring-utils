package model

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	EditorModel struct {
		err     error
		records table.Model
	}
)

func initEditor(records [][]string) EditorModel {
	rows := []table.Row{}
	for i, r := range records {
		r = append([]string{fmt.Sprintf("%d", i+1)}, r...)
		rows = append(rows, table.Row(r))
	}

	cols := []table.Column{
		{Title: "Row", Width: 3},
		{Title: "Name", Width: 25},
		{Title: "Email", Width: 25},
	}

	t := table.New(
		table.WithColumns(cols),
		table.WithRows(rows),
		table.WithStyles(table.Styles{
			Header:   lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("210")),
			Selected: lipgloss.NewStyle().Background(lipgloss.Color("236")).Foreground(lipgloss.Color("15")),
		}),
	)

	t.Focus()

	return EditorModel{
		records: t,
	}
}

func (e EditorModel) Init() tea.Cmd {
	return nil
}

func (e EditorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	e.records, cmd = e.records.Update(msg)
	return e, cmd
}

func (e EditorModel) View() string {
	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63")).
		Render(e.records.View())
}
