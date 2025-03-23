// Copyright © 2025 Duane Matthew P. Chan

package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	email "github.com/duanechan/monitoring-utils/email/internal"
	"github.com/duanechan/monitoring-utils/email/internal/model"
)

func main() {
	rName := flag.String("name", "", "the name of the recipient")
	rEmail := flag.String("email", "", "the email of the recipient")
	flag.Parse()

	config, err := email.LoadConfig()
	if err != nil {
		fmt.Printf("error loading config: %s\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(model.InitializeModel(*rName, *rEmail, config))
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
