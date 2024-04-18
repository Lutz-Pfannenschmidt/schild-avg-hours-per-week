package csv

import (
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/xuri/excelize/v2"
)

func ResultToCSV(in map[string][2]float64) string {
	var result = "Name,1. Halbjahr,2. Halbjahr\n"
	for name, values := range in {
		result += fmt.Sprintf("%s,%.2f,%.2f\n", name, values[0], values[1])
	}
	return strings.ReplaceAll(result, "-1.00", "No Data")
}

func WriteToFile(path string, in map[string][2]float64) error {
	content := ResultToCSV(in)
	return os.WriteFile(path, []byte(content), fs.ModePerm)
}

func ReadXLSXFileToCSV(path string) (*[]string, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("no sheets found")
	}

	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return nil, err
	}

	result := []string{}

	for _, row := range rows {
		result = append(result, strings.Join(row, ","))
	}

	return &result, nil
}

func ReadCSVFile(path string) (*[]string, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(file), "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}
	return &lines, nil
}

func ReadAnyFileToCSV(path string) (*[]string, error) {
	if strings.HasSuffix(path, ".csv") {
		return ReadCSVFile(path)
	} else if strings.HasSuffix(path, ".xlsx") {
		return ReadXLSXFileToCSV(path)
	} else {
		return nil, fmt.Errorf("unsupported file type")
	}
}
