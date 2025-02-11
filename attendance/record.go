package main

import (
	"fmt"
	"strings"
)

type Record struct {
	In    string
	Out   string
	Date  string
	Email string
}

type Records []Record

func NewRecord(date, email, in, out string) Record {
	return Record{Date: date, Email: email, In: in, Out: out}
}

func (r Record) String() string {
	return fmt.Sprintf("%s %s %s %s\n", r.Date, r.Email, r.In, r.Out)
}

func (r Records) ContainsValue(v string) bool {
	for _, record := range r {
		if strings.Contains(record.String(), v) {
			return true
		}
	}
	return false
}

func (r Records) Display() {
	fmt.Println(" ┌────────────┬───────────────────────────────────────────┬──────────┬──────────┐")
	fmt.Println(" │ Date       │ Email                                     │ Time in  │ Time out │")
	fmt.Println(" ├────────────┼───────────────────────────────────────────┼──────────┼──────────┤")

	for _, record := range r {
		fmt.Printf("│ %-11s │ %-41s │ %-8s │ %-8s │\n",
			record.Date,
			record.Email,
			record.In,
			record.Out,
		)
	}

	fmt.Println(" └────────────┴───────────────────────────────────────────┴──────────┴──────────┘")
}
