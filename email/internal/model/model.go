package model

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type EmailModel struct {
	helpShown bool
	err       error

	parser ParserModel
	sender SenderModel
}

func InitializeModel() EmailModel {
	return EmailModel{
		parser: initParser(),
		sender: initSender(),
	}
}

// func (e *EmailModel) readInput() {
// 	records, err := email.ParseData(e.textInput.Value())
// 	if err != nil {
// 		e.err = err
// 		return
// 	}

// 	for i, record := range records {
// 		e.textArea.InsertString(fmt.Sprintf("%s,%s", record[0], record[1]))
// 		if i != len(records)-1 {
// 			e.textArea.InsertString("\n")
// 		}
// 	}

// 	e.result = email.ValidateRecords(records)

// 	e.textInput.Reset()
// 	e.textInput.Blur()
// }

func (e EmailModel) Init() tea.Cmd {
	return tea.Batch(tea.ClearScreen, textarea.Blink)
}

func (e EmailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "?":
			e.helpShown = !e.helpShown
			return e, nil
		case "esc":
			fallthrough
		case "ctrl+c":
			return e, tea.Quit
		case "enter":
			e.err = nil
		}
	case parserError:
		e.err = msg.err
		return e, nil
	}

	if e.parser.result != nil {

	} else {
		var parser tea.Model
		parser, cmd = e.parser.Update(msg)
		e.parser = parser.(ParserModel)
	}
	return e, cmd
}

func (e EmailModel) View() string {
	sections := []string{}

	sections = append(sections, e.headerView())

	if e.parser.result == nil {
		sections = append(sections, e.parser.View())
	}

	if e.err != nil {
		sections = append(sections, e.errorView())
	} else {
		sections = append(sections, "\n\n")
	}

	sections = append(sections, e.helpView())

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}
