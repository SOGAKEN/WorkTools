package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

// Record represents a row in the CSV file.
type Record struct {
	Switch string
	Inum   string
}

// findFiles finds the files that start with specified prefixes.
func findFiles(prefix1, prefix2 string) (string, string, error) {
	entries, err := os.ReadDir(".")
	if err != nil {
		return "", "", err
	}

	var file1, file2 string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if strings.HasPrefix(name, prefix1) && file1 == "" {
			file1 = name
		} else if strings.HasPrefix(name, prefix2) && file2 == "" {
			file2 = name
		}
	}

	if file1 == "" || file2 == "" {
		return "", "", fmt.Errorf("could not find matching files")
	}

	return file1, file2, nil
}

// getColumnIndex finds the index of a column in the header.
func getColumnIndex(header []string, columnName string) (int, error) {
	for i, h := range header {
		if h == columnName {
			return i, nil
		}
	}
	return -1, fmt.Errorf("column %s not found", columnName)
}

// ReadCSV reads a CSV file and returns a map of records and the header.
func ReadCSV(filename string) (map[string]Record, []string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, nil, err
	}

	if len(records) < 2 {
		return nil, nil, fmt.Errorf("file %s does not contain enough data", filename)
	}

	header := records[0]

	switchIndex, err := getColumnIndex(header, "switch")
	if err != nil {
		return nil, nil, err
	}
	inumIndex, err := getColumnIndex(header, "inum")
	if err != nil {
		return nil, nil, err
	}

	recordMap := make(map[string]Record)
	for _, record := range records[1:] {
		key := record[switchIndex] + "-" + record[inumIndex]
		recordMap[key] = Record{Switch: record[switchIndex], Inum: record[inumIndex]}
	}

	return recordMap, header, nil
}

// writeCSV writes records to a CSV file, including the header.
func writeCSV(filename string, records []Record, header []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the header
	if err := writer.Write(header); err != nil {
		return err
	}

	for _, record := range records {
		if err := writer.Write([]string{record.Switch, record.Inum}); err != nil {
			return err
		}
	}

	return nil
}

// CompareAndWrite compares two CSV files and writes the results to separate files.
func CompareAndWrite(fileA, fileB string) error {
	recordsA, headerA, err := ReadCSV(fileA)
	if err != nil {
		return err
	}

	recordsB, headerB, err := ReadCSV(fileB)
	if err != nil {
		return err
	}

	aOnly, bOnly, common := make([]Record, 0), make([]Record, 0), make([]Record, 0)

	for key, record := range recordsA {
		if _, exists := recordsB[key]; !exists {
			aOnly = append(aOnly, record)
		} else {
			common = append(common, record)
		}
	}

	for key, record := range recordsB {
		if _, exists := recordsA[key]; !exists {
			bOnly = append(bOnly, record)
		}
	}

	// Define output file names based on input file names
	prefixA := strings.TrimSuffix(fileA, ".csv")
	prefixB := strings.TrimSuffix(fileB, ".csv")

	// Write results to files
	writeCSV(prefixA+"_a_only.csv", aOnly, headerA)
	writeCSV(prefixB+"_b_only.csv", bOnly, headerB)
	writeCSV(prefixA+"_common.csv", common, headerA)

	return nil
}

func main() {
	osakaFile, tokyoFile, err := findFiles("osaka_", "tokyo_")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if err := CompareAndWrite(osakaFile, tokyoFile); err != nil {
		fmt.Println("Error:", err)
	}
}
