package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

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

// CompareFiles compares two CSV files and writes the results to separate files.
// osaka_ ファイルを基準にして、tokyo_ ファイルを比較します。
func CompareFiles(fileOsaka, fileTokyo string) error {
	recordsOsaka, headerOsaka, err := ReadCSV(fileOsaka)
	if err != nil {
		return err
	}

	recordsTokyo, _, err := ReadCSV(fileTokyo) // headerTokyo は使用しないので、_ で無視
	if err != nil {
		return err
	}

	inTokyo, notInTokyo := make([][]string, 0), make([][]string, 0)

	// osaka_ ファイルの各レコードに対して、tokyo_ ファイルに含まれているかを確認
	for key, record := range recordsOsaka {
		if _, exists := recordsTokyo[key]; exists {
			inTokyo = append(inTokyo, record)
		} else {
			notInTokyo = append(notInTokyo, record)
		}
	}

	// 出力ファイル名を定義
	prefixOsaka := strings.TrimSuffix(fileOsaka, ".csv")

	// 結果をファイルに書き出す
	writeCSV(prefixOsaka+"_in_tokyo.csv", inTokyo, headerOsaka)
	writeCSV(prefixOsaka+"_not_in_tokyo.csv", notInTokyo, headerOsaka)

	return nil
}

func main() {
	osakaFile, tokyoFile, err := findFiles("osaka_", "tokyo_")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if err := CompareFiles(osakaFile, tokyoFile); err != nil {
		fmt.Println("Error:", err)
	}
}
