package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

func main() {
	// 各CSVファイルからinum値を読み込む
	aInums, err := readCsv("a.csv")
	if err != nil {
		panic(err)
	}
	bInums, err := readCsv("b.csv")
	if err != nil {
		panic(err)
	}

	// リストを比較して結果を得る
	aUnique, bUnique, common := compareLists(aInums, bInums)

	// 結果をCSVファイルに出力する
	writeCsv("a_unique.csv", aUnique)
	writeCsv("b_unique.csv", bUnique)
	writeCsv("common.csv", common)

	fmt.Println("CSV files have been written successfully.")
}

// readCsv はCSVファイルを読み込み、inum値のスライスを返す
func readCsv(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var inums []string
	for _, record := range records {
		inums = append(inums, record[0]) // 'inum'が最初の列であると仮定
	}

	return inums, nil
}

// compareLists は2つのスライスを比較し、ユニークな要素と共通要素を返す
func compareLists(a, b []string) (aUnique, bUnique, common []string) {
	mA := make(map[string]bool)
	mB := make(map[string]bool)

	for _, item := range a {
		mA[item] = true
	}
	for _, item := range b {
		mB[item] = true
		if _, found := mA[item]; found {
			common = append(common, item)
		} else {
			bUnique = append(bUnique, item)
		}
	}
	for _, item := range a {
		if _, found := mB[item]; !found {
			aUnique = append(aUnique, item)
		}
	}

	return
}

// writeCsv は指定されたファイル名でスライスをCSVファイルに書き込む
func writeCsv(filename string, data []string) {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, value := range data {
		if err := writer.Write([]string{value}); err != nil {
			panic("Error writing record to csv:" + err.Error())
		}
	}
}
