package main

import (
	"log"
	"os"
	"strings"
)

func GetRecords(file string) Records {
	bytes, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf("error: failed to read data.csv: %s", err)
	}

	data := strings.Split(string(bytes), "\n")
	records := Records{}
	for _, d := range data {
		recordData := strings.Split(d, ",")
		if len(recordData) < 4 || recordData[3] == "" {
			recordData = append(recordData, "00:00:00")
		}
		record := NewRecord(
			recordData[0],
			recordData[1],
			recordData[2],
			recordData[3],
		)

		records = append(records, record)
	}

	return Search(records, Today)
}

func main() {
	records := GetRecords("data.csv")
	records.Display()
}
