// Copyright © 2025 Duane Matthew P. Chan

package model

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	headerStyle lipgloss.Style = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#6796bf"))
)

type CredentialsModel struct {
	Parser parserModel
	// EmailSender senderModel
}

func InitializeModel(path string) CredentialsModel {
	return CredentialsModel{
		Parser: initParser(),
	}
}

func (m CredentialsModel) Init() tea.Cmd {
	m.Parser.Init()
	return tea.ClearScreen
}

func (m CredentialsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc, tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEnter:
			return m, func() tea.Msg {
				input := m.Parser.textInput.Value()
				return parseStart{input}
			}
		}

	}

	var parser tea.Model
	parser, cmd = m.Parser.Update(msg)
	m.Parser = parser.(parserModel)
	return m, cmd
}

func (m CredentialsModel) View() string {
	return fmt.Sprintf("\n\n%s\n\n\n\n\n%s\n\n",
		headerStyle.Render(` _____              _            _   _       _       _   _      _                 
/  __ \            | |          | | (_)     | |     | | | |    | |                
| /  \/_ __ ___  __| | ___ _ __ | |_ _  __ _| |___  | |_| | ___| |_ __   ___ _ __ 
| |   | '__/ _ \/ _`+"`"+` |/ _ | '_ \| __| |/ _`+"`"+` | / __| |  _  |/ _ | | '_ \ / _ | '__|
| \__/| | |  __| (_| |  __| | | | |_| | (_| | \__ \ | | | |  __| | |_) |  __| |   
 \____|_|  \___|\__,_|\___|_| |_|\__|_|\__,_|_|___/ \_| |_/\___|_| .__/ \___|_|   
                                                                 | |              
                                                                 |_|     `),
		m.Parser.View(),
	)
}
