// Copyright © 2025 Duane Matthew P. Chan

package credentials

import tea "github.com/charmbracelet/bubbletea"

type CredentialsModel struct {
}

func InitializeModel() CredentialsModel {
	return CredentialsModel{}
}

func (m CredentialsModel) Init() tea.Cmd {
	return nil
}

func (m CredentialsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return nil, nil
}

func (m CredentialsModel) View() string {
	return ""
}
