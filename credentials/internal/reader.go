package credentials

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/xuri/excelize/v2"
)

type ReadAll func() ([][]string, error)

// Identifies the file type and returns a ReadAll function and an error.
func NewReader(filepath string) (ReadAll, error) {
	switch {
	// CSV
	case strings.HasSuffix(filepath, ".csv"):
		return func() ([][]string, error) {
			return ReadCSV(filepath)
		}, nil

	// XLSX
	case strings.HasSuffix(filepath, ".xlsx"):
		return func() ([][]string, error) {
			return ReadXLSX(filepath)
		}, nil

	// Default
	default:
		return nil, fmt.Errorf("file not supported")
	}
}

func ReadCSV(filepath string) ([][]string, error) {
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

func ReadXLSX(filepath string) ([][]string, error) {
	file, err := excelize.OpenFile(filepath)
	if err != nil {
		return [][]string{}, err
	}
	defer file.Close()

	records, err := file.GetRows(file.GetSheetName(0))
	if err != nil {
		return [][]string{}, err
	}

	return records, nil
}
