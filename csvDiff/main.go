package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	folderPath := "." // Set the folder path here

	csvFiles, err := findCSVFiles(folderPath)
	if err != nil {
		fmt.Println("Error finding CSV files:", err)
		return
	}

	if len(csvFiles) < 2 {
		fmt.Println("Need at least two CSV files for comparison")
		return
	}

	// Read and store data from all CSV files
	allRecords := make(map[string][][]string)
	allHeaders := make(map[string][]string)
	for _, file := range csvFiles {
		records, header, err := readCSV(file)
		if err != nil {
			fmt.Println("Error reading file:", file, err)
			continue
		}
		allRecords[file] = records
		allHeaders[file] = header
	}

	// Compare and classify records from the first two CSV files
	file1 := csvFiles[0]
	file2 := csvFiles[1]

	switchIndex1 := findColumnIndex(allHeaders[file1], "switch")
	switchIndex2 := findColumnIndex(allHeaders[file2], "switch")

	if switchIndex1 == -1 || switchIndex2 == -1 {
		fmt.Println("Error: 'switch' column not found in one or both files")
		return
	}

	switches1 := getSwitchValues(allRecords[file1], switchIndex1)
	switches2 := getSwitchValues(allRecords[file2], switchIndex2)

	uniqueFile1, uniqueFile2, common := classifyRecords(switches1, switches2, allRecords[file1], allRecords[file2], switchIndex1, switchIndex2)

	// Save results to CSV files with headers
	saveCSV("unique_"+filepath.Base(file1), uniqueFile1, allHeaders[file1])
	saveCSV("unique_"+filepath.Base(file2), uniqueFile2, allHeaders[file2])
	saveCSV("common_records.csv", common, allHeaders[file1]) // Assuming both files have the same header structure
}

// findCSVFiles finds and returns all CSV files in the given directory
func findCSVFiles(folderPath string) ([]string, error) {
	var files []string
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".csv" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func readCSV(filename string) ([][]string, []string, error) {
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

	header := records[0] // First row is the header
	data := records[1:]  // Remaining rows are data
	return data, header, nil
}

func findColumnIndex(header []string, columnName string) int {
	for index, column := range header {
		if column == columnName {
			return index
		}
	}
	return -1 // Column not found
}

func getSwitchValues(records [][]string, switchIndex int) map[string]bool {
	switches := make(map[string]bool)
	for _, record := range records {
		if len(record) > switchIndex {
			switches[record[switchIndex]] = true
		}
	}
	return switches
}

// classifyRecords function modified to include duplicates in common records
func classifyRecords(tokyoSwitches, osakaSwitches map[string]bool, tokyoRecords, osakaRecords [][]string, switchIndexTokyo, switchIndexOsaka int) ([][]string, [][]string, [][]string) {
	uniqueTokyo := [][]string{}
	uniqueOsaka := [][]string{}
	common := [][]string{}

	for _, record := range tokyoRecords {
		if len(record) > switchIndexTokyo {
			_, inOsaka := osakaSwitches[record[switchIndexTokyo]]
			if !inOsaka {
				uniqueTokyo = append(uniqueTokyo, record)
			} else {
				common = append(common, record) // Add to common if found in Osaka
			}
		}
	}

	for _, record := range osakaRecords {
		if len(record) > switchIndexOsaka {
			_, inTokyo := tokyoSwitches[record[switchIndexOsaka]]
			if !inTokyo {
				uniqueOsaka = append(uniqueOsaka, record)
			} else {
				common = append(common, record) // Add to common if found in Tokyo
			}
		}
	}

	return uniqueTokyo, uniqueOsaka, common
}

func saveCSV(filename string, records [][]string, header []string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write(header); err != nil { // Write the header first
		fmt.Println("Error writing header to file:", err)
		return
	}

	for _, record := range records {
		if err := writer.Write(record); err != nil {
			fmt.Println("Error writing record to file:", err)
			return
		}
	}
}
