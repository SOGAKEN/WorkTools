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

	for _, aFile := range aFiles {
		aData, aHeader := readCsv(aFile)

		for _, bFile := range bFiles {
			bData, _ := readCsv(bFile)

			aOnly, bOnly, both := compareCsv(aData, bData)

			writeCsv("only_"+aFile, aOnly, aHeader)
			writeCsv("only_"+bFile, bOnly, aHeader)
			writeCsv("both_"+aFile+"_"+bFile, both, aHeader)
		}
	}
}

func readCsv(filename string) ([][]string, []string) {
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
	return records[1:], header // ヘッダーとデータ行を返す
}

func compareCsv(aData, bData [][]string) ([][]string, [][]string, [][]string) {
	aOnly := [][]string{}
	bOnly := [][]string{}
	both := [][]string{}

	bMap := make(map[string]bool)
	for _, b := range bData {
		bMap[strings.Join(b, ",")] = true
	}

	for _, a := range aData {
		aStr := strings.Join(a, ",")
		if _, found := bMap[aStr]; found {
			both = append(both, a)
			delete(bMap, aStr) // 同じ行を再度検討しないように削除
		} else {
			aOnly = append(aOnly, a)
		}
	}

	for _, b := range bData {
		bStr := strings.Join(b, ",")
		if _, found := bMap[bStr]; found {
			bOnly = append(bOnly, b)
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
