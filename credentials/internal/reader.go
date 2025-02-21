// Copyright © 2025 Duane Matthew P. Chan

package credentials

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

// Parses raw (CSV file) data and returns a slice of recipients.
func GetRecipients(filepath string) ([]User, int, int, error) {
	if !strings.HasSuffix(filepath, ".csv") {
		return []User{}, -1, -1, fmt.Errorf("file is not CSV")
	}

	file, err := os.Open(filepath)
	if err != nil {
		return []User{}, -1, -1, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return []User{}, -1, -1, err
	}

	emailMap, invalidEmails, numOfDupes := ValidateRecords(records)

	if len(invalidEmails) > 0 {
		ContinuePrompt(invalidEmails)
	}

	recipients := []User{}
	for k, v := range emailMap {
		recipients = append(recipients, User{Name: v, Email: k})
	}

	return recipients, len(invalidEmails), numOfDupes, nil
}

func ContinuePrompt(invalidEmails map[int]string) {
	fmt.Println()
	color.HiYellow("There is/are %d invalid email/s in the file:\n", len(invalidEmails))
	for k, v := range invalidEmails {
		fmt.Printf("-> Record (row) %d: %s\n", k, v)
	}

	fmt.Printf(
		"\nAre you sure you want to continue? Press Enter to %s or CTRL+C to %s.",
		color.New(color.FgHiYellow).Sprintf("continue"),
		color.New(color.FgHiGreen).Sprintf("cancel"),
	)
	fmt.Scanln()
	fmt.Println()
}

func ValidateRecords(records [][]string) (map[string]string, map[int]string, int) {
	recipients := map[string]string{}
	invalid := map[int]string{}
	duplicates := 0

	for i, r := range records {
		name := strings.TrimSpace(strings.ReplaceAll(r[0], "\r", ""))
		email := strings.TrimSpace(strings.ReplaceAll(r[1], "\r", ""))

		if !IsValidEmail(email) {
			invalid[i+1] = fmt.Sprintf("Invalid email address (%s).", email)
		}

		if _, exists := recipients[email]; exists {
			dupeIdx := 1

			for k := range recipients {
				if k == email {
					break
				}
				dupeIdx++
			}

			invalid[i+1] = fmt.Sprintf("Duplicate email. Exact match at record %d (%s).", dupeIdx, email)
			duplicates++
		}

		recipients[email] = name
	}

	return recipients, invalid, duplicates
}
