package model

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (e EmailModel) headerView() string {
	return Header.Render(`
	______     __    __     ______     __     __            __  __     ______     __         ______   ______     ______    
   /\  ___\   /\ "-./  \   /\  __ \   /\ \   /\ \          /\ \_\ \   /\  ___\   /\ \       /\  == \ /\  ___\   /\  == \   
   \ \  __\   \ \ \-./\ \  \ \  __ \  \ \ \  \ \ \____     \ \  __ \  \ \  __\   \ \ \____  \ \  _-/ \ \  __\   \ \  __<   
	\ \_____\  \ \_\ \ \_\  \ \_\ \_\  \ \_\  \ \_____\     \ \_\ \_\  \ \_____\  \ \_____\  \ \_\    \ \_____\  \ \_\ \_\ 
	 \/_____/   \/_/  \/_/   \/_/\/_/   \/_/   \/_____/      \/_/\/_/   \/_____/   \/_____/   \/_/     \/_____/   \/_/ /_/ 
																														   
	   `)
}

func (e EmailModel) helpView() string {
	if !e.helpShown {
		return Help.Render(lipgloss.NewStyle().Padding(1, 2).Render("? / toggle help"))
	}

	commands := lipgloss.JoinHorizontal(
		lipgloss.Center,
		lipgloss.NewStyle().Padding(1, 2).Render("? / toggle help"),
		lipgloss.NewStyle().Padding(1, 2).Render("ctrl+c, esc / exit program"),
	)

	return Help.Render(commands)
}

func (e EmailModel) errorView() string {
	if e.err != nil {
		return Error.Render(fmt.Sprint("Error: ", e.err))
	}
	return ""
}
