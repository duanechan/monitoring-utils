package model

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	email "github.com/duanechan/monitoring-utils/email/internal"
)

type (
	// ParserModel
	ParserModel struct {
		enabled   bool
		result    *email.ParseResult
		textInput textinput.Model
	}

	// ParserError
	parserError struct{ err error }
)

func initParser() ParserModel {
	ti := textinput.New()
	ti.Prompt = "Input the filepath: "
	ti.PromptStyle = lipgloss.NewStyle().
		AlignHorizontal(lipgloss.Center).
		Bold(true)
	ti.Placeholder = "C:/path/to/file"
	ti.Width = 150
	ti.CharLimit = 150
	ti.Focus()

	return ParserModel{
		enabled:   true,
		result:    nil,
		textInput: ti,
	}
}

func (p ParserModel) Init() tea.Cmd {
	return textinput.Blink
}

func (p ParserModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return p, tea.Quit
		case "enter":
			if err := p.readInput(); err != nil {
				return p, func() tea.Msg {
					return parserError{err}
				}
			}
		}
	}

	if p.enabled {
		p.textInput, cmd = p.textInput.Update(msg)
	}
	return p, cmd
}

func (p ParserModel) View() string {
	return p.textInput.View()
}

func (p *ParserModel) readInput() error {
	input := strings.ReplaceAll(p.textInput.Value(), "\"", "")
	records, err := email.ParseData(input)
	if err != nil {
		return err
	}

	p.result = email.ValidateRecords(records)

	p.textInput.Reset()
	p.textInput.Blur()

	return nil
}
