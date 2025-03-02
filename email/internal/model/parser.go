package model

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	email "github.com/duanechan/monitoring-utils/email/internal"
)

type (
	ParserModel struct {
		enabled   bool
		err       error
		result    email.ParseResult
		textInput textinput.Model
	}

	parserMessage struct {
		err    error
		result email.ParseResult
	}
	readInputMessage struct{}
)

func initParser() ParserModel {
	ti := textinput.New()
	ti.Prompt = "Input the filepath: "
	ti.PromptStyle = lipgloss.NewStyle().
		AlignHorizontal(lipgloss.Center).
		Bold(true)
	ti.Placeholder = "C:/path/to/file"
	ti.Width = 100
	ti.CharLimit = 150
	ti.Focus()

	return ParserModel{
		enabled:   true,
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
		case "enter":
			return p, p.readInput
		}

	case readInputMessage:
		p.result, p.err = p.parse()
		return p, p.parserSuccess
	}

	if p.enabled {
		p.textInput, cmd = p.textInput.Update(msg)
	}
	return p, cmd
}

func (p ParserModel) View() string {
	return lipgloss.NewStyle().
		Padding(1, 2).
		Render(p.textInput.View())
}

func (p ParserModel) readInput() tea.Msg {
	return readInputMessage{}
}

func (p ParserModel) parserSuccess() tea.Msg {
	return parserMessage{
		err:    p.err,
		result: p.result,
	}
}

func (p *ParserModel) parse() (email.ParseResult, error) {
	input := strings.ReplaceAll(p.textInput.Value(), "\"", "")
	records, err := email.ParseData(input)
	if err != nil {
		return email.ParseResult{}, err
	}

	result := email.ValidateRecords(records)

	p.textInput.Reset()
	p.textInput.Blur()

	return result, nil
}
