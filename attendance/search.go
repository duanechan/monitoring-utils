package main

import (
	"strings"
	"time"
)

func Today(s string) bool {
	return strings.Contains(s, time.Now().Format("1/2/2006"))
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
