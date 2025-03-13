// Copyright © 2025 Duane Matthew P. Chan

package email

import (
	"fmt"
	"strings"
)

type ParseResult struct {
	Invalids   int
	Duplicates int
	Raw        [][]string
	Recipients []User
	BadEmails  map[int]string
}

func (p ParseResult) IsEmpty() bool {
	return len(p.Recipients) == 0 &&
		len(p.Raw) == 0 &&
		p.Invalids == 0 &&
		p.Duplicates == 0 &&
		len(p.BadEmails) == 0
}

// Parses raw (CSV file) data and returns a slice of recipients.
func ParseData(filepath string) ([][]string, error) {
	if filepath == "" {
		return [][]string{}, fmt.Errorf("no filepath provided")
	}

	ReadAll, err := NewReader(filepath)
	if err != nil {
		return [][]string{}, err
	}

	records, err := ReadAll()
	if err != nil {
		return [][]string{}, err
	}

	return records, nil
}

func ValidateRecords(records [][]string) ParseResult {
	recipientMap := map[string]int{}
	result := ParseResult{Raw: records, BadEmails: make(map[int]string)}

	for i, r := range records {
		name := strings.TrimSpace(strings.ReplaceAll(r[0], "\r", ""))
		email := strings.TrimSpace(strings.ReplaceAll(r[1], "\r", ""))

		if !IsValidEmail(email) {
			result.Invalids++
			result.BadEmails[i+1] = fmt.Sprintf("Invalid email address at row %d (%s).", i+1, email)
			continue
		}
		if dupeIdx, exists := recipientMap[email]; exists {
			result.Duplicates++
			result.BadEmails[i+1] = fmt.Sprintf("Duplicate email at row %d. Exact match at record %d (%s).", i+1, dupeIdx+1, email)
			continue
		} else {
			recipientMap[email] = i
		}

		result.Recipients = append(result.Recipients, User{Name: name, Email: email})
	}

	return result
}
