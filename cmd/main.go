// Copyright © 2025 Duane Matthew P. Chan

package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/duanechan/monitoring-utils/email/internal/model"
)

func main() {
	p := tea.NewProgram(model.InitializeModel())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
