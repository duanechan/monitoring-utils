// Copyright © 2025 Duane Matthew P. Chan

package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/duanechan/monitoring-utils/email/internal/model"
)

func main() {
	name := flag.String("name", "", "the name of the recipient")
	email := flag.String("email", "", "the email of the recipient")
	flag.Parse()

	p := tea.NewProgram(model.InitializeModel(*name, *email))
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
