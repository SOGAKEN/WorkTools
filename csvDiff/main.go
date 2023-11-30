package main

import (
	"encoding/csv"
	"os"
	"strings"
)

func main() {
	files, err := os.ReadDir(".")
	if err != nil {
		panic(err)
	}

	var aFiles, bFiles []string

	for _, file := range files {
		if strings.Contains(file.Name(), "A") {
			aFiles = append(aFiles, file.Name())
		}
		if strings.Contains(file.Name(), "B") {
			bFiles = append(bFiles, file.Name())
		}
	}

	var aHeader []string

	for _, aFile := range aFiles {
		aData, header := readCsv(aFile)
		aHeader = header // A.csvのヘッダーを取得

		for _, bFile := range bFiles {
			bData, _ := readCsv(bFile) // B.csvのヘッダーは無視

			aOnly, bOnly, both := compareCsv(aData, bData)

			writeCsv("only_"+aFile, aOnly, aHeader)
			writeCsv("only_"+bFile, bOnly, aHeader)
			writeCsv("both_"+aFile+"_"+bFile, both, aHeader)
		}
	}
}

func readCsv(filename string) (map[string][]string, []string) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	if len(records) == 0 {
		return nil, nil
	}

	header := records[0]
	inumIndex := -1
	for i, col := range header {
		if col == "inum" {
			inumIndex = i
			break
		}
	}

	if inumIndex == -1 {
		panic("inum column not found")
	}

	data := make(map[string][]string)
	for _, record := range records[1:] { // 最初のヘッダー行をスキップ
		inum := record[inumIndex]
		data[inum] = record
	}

	return data, header
}

func compareCsv(aData, bData map[string][]string) ([][]string, [][]string, [][]string) {
	aOnly := [][]string{}
	bOnly := [][]string{}
	both := [][]string{}

	for inum, record := range aData {
		if _, ok := bData[inum]; ok {
			both = append(both, record)
		} else {
			aOnly = append(aOnly, record)
		}
	}

	for inum, record := range bData {
		if _, ok := aData[inum]; !ok {
			bOnly = append(bOnly, record)
		}
	}

	return aOnly, bOnly, both
}

func writeCsv(filename string, data [][]string, header []string) {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	if len(header) > 0 {
		if err := writer.Write(header); err != nil {
			panic(err)
		}
	}

	for _, record := range data {
		if err := writer.Write(record); err != nil {
			panic(err)
		}
	}

	writer.Flush()
}
