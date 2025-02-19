package main

import (
	"encoding/csv"
	"log"
	"os"
)

func GetRecords(filepath string) (Records, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return Records{}, err
	}

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return Records{}, err
	}

	var employees Records

	for _, r := range records {
		employees = append(employees, Record{
			In:    r[2],
			Out:   r[3],
			Date:  r[0],
			Email: r[1],
		})
	}

	return Search(employees, Today), nil
}

func main() {
	records, err := GetRecords("data.csv")
	if err != nil {
		log.Fatalf("error: failed to get records: %s", err)
	}
	records.Display()
}
