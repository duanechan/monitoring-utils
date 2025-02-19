package main

import (
	"strings"
)

func Today(s string) bool {
	return strings.Contains(s, "2/18/2025")
}

func Search(records Records, predicate func(s string) bool) Records {
	filtered := Records{}

	for _, r := range records {
		if predicate(r.String()) {
			filtered = append(filtered, r)
		}
	}

	return filtered
}
