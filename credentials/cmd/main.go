// Copyright © 2025 Duane Matthew P. Chan

package main

import (
	"flag"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/duanechan/monitoring-utils/credentials/internal/model"
)

func main() {
	path := flag.String("path", "", "the file path to the list of recipients")
	flag.Parse()

	p := tea.NewProgram(model.InitializeModel(*path))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
