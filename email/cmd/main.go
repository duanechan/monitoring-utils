// Copyright © 2025 Duane Matthew P. Chan

package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/duanechan/monitoring-utils/email/internal/model"
)

func main() {
	p := tea.NewProgram(model.InitializeModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
