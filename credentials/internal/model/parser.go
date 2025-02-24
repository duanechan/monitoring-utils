package model

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	credentials "github.com/duanechan/monitoring-utils/credentials/internal"
)

type (
	// Parser model
	parserModel struct {
		invalidRecords   int
		duplicateRecords int
		errorLog         string
		recipients       []credentials.User
		textInput        textinput.Model
	}

	parseResult struct {
		invalids   int
		duplicates int
		errorLog   string
		recipients []credentials.User
	}

	// Parser messages
	// parseStart message signals the parser to read their initialized filepath
	parseStart   struct{ string }
	parseSuccess struct{ parseResult }
	parseError   struct{ error }

	// Parser errors

)

func initParser() parserModel {
	ti := textinput.New()
	ti.Placeholder = "C:/path/to/recipients.csv"
	ti.Focus()
	ti.CharLimit = 200
	ti.Width = 200

	return parserModel{
		textInput: ti,
	}
}

func (p parserModel) Init() tea.Cmd {
	return textinput.Blink
}

func (p parserModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case parseStart:
		return p, func() tea.Msg {
			records, err := parseData(msg.string)
			if err != nil {
				return parseError{err}
			}

			result := validateRecords(records)

			return parseSuccess{result}
		}

	case parseSuccess:
		p.recipients = msg.recipients
		p.invalidRecords = msg.invalids
		p.duplicateRecords = msg.duplicates
		p.textInput.Reset()
		p.textInput.Blur()
		return p, nil
	}

	p.textInput, cmd = p.textInput.Update(msg)
	return p, cmd
}

func (p parserModel) View() string {
	if len(p.recipients) > 0 {
		return fmt.Sprintf("%s\n%s", p.errorLog, p.recipients)
	}
	return p.textInput.View()
}

func parseData(filepath string) ([][]string, error) {
	if !strings.HasSuffix(filepath, ".csv") {
		return [][]string{}, fmt.Errorf("file type is not .csv")
	}

	file, err := os.Open(filepath)
	if err != nil {
		return [][]string{}, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return [][]string{}, err
	}

	return records, nil
}

func validateRecords(records [][]string) parseResult {
	recipientMap := map[string]int{}
	result := parseResult{}

	for i, r := range records {
		name := strings.TrimSpace(strings.ReplaceAll(r[0], "\r", ""))
		email := strings.TrimSpace(strings.ReplaceAll(r[1], "\r", ""))

		if !credentials.IsValidEmail(email) {
			result.invalids++
			result.errorLog += fmt.Sprintf("Invalid email address (%s).\n", email)
			continue
		}

		if dupeIdx, exists := recipientMap[email]; exists {
			result.duplicates++
			result.errorLog += fmt.Sprintf("Duplicate email. Exact match at record %d (%s).\n", dupeIdx+1, email)
			continue
		} else {
			recipientMap[email] = i
		}

		result.recipients = append(result.recipients, credentials.User{Name: name, Email: email})
	}

	return result
}
