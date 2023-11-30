package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

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

// trimDecimal trims the decimal part of the inum value.
func trimDecimal(inum string) string {
	if dot := strings.Index(inum, "."); dot != -1 {
		return inum[:dot]
	}
	return inum
}

// ReadCSV reads a CSV file and returns a map of the records (entire row) and the header.
func ReadCSV(filename string) (map[string][]string, []string, error) {
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

	recordMap := make(map[string][]string)
	for _, record := range records[1:] {
		inum := trimDecimal(record[inumIndex])
		key := record[switchIndex] + "-" + inum
		recordMap[key] = record
	}

	return recordMap, header, nil
}

// writeCSV writes records to a CSV file, including the header.
func writeCSV(filename string, records [][]string, header []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write(header); err != nil {
		return err
	}

	for _, record := range records {
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

// CompareAndWrite compares two CSV files and writes the results to separate files.
func CompareFiles(fileTokyo, fileOsaka string) error {
	recordsTokyo, headerTokyo, err := ReadCSV(fileTokyo)
	if err != nil {
		return err
	}

	recordsOsaka, _, err := ReadCSV(fileOsaka) // headerOsaka は使用しないので、_ で無視
	if err != nil {
		return err
	}

	inOsaka, notInOsaka := make([][]string, 0), make([][]string, 0)

	// tokyo_ ファイルの各レコードに対して、osaka_ ファイルに含まれているかを確認
	for key, record := range recordsTokyo {
		if _, exists := recordsOsaka[key]; exists {
			inOsaka = append(inOsaka, record)
		} else {
			notInOsaka = append(notInOsaka, record)
		}
	}

	// 出力ファイル名を定義
	prefixTokyo := strings.TrimSuffix(fileTokyo, ".csv")

	// 結果をファイルに書き出す
	writeCSV(prefixTokyo+"_in_osaka.csv", inOsaka, headerTokyo)
	writeCSV(prefixTokyo+"_not_in_osaka.csv", notInOsaka, headerTokyo)

	return nil
}

func main() {
	tokyoFile, osakaFile, err := findFiles("tokyo_", "osaka_")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if err := CompareFiles(tokyoFile, osakaFile); err != nil {
		fmt.Println("Error:", err)
	}
}
