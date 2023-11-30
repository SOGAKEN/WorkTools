package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

func main() {
	// 各CSVファイルからデータを読み込む
	aData, header, err := readCsv("a.csv")
	if err != nil {
		panic(err)
	}
	bData, _, err := readCsv("b.csv")
	if err != nil {
		panic(err)
	}

	// リストを比較して結果を得る
	aUnique, bUnique, common := compareLists(aData, bData)

	// 結果をCSVファイルに出力する
	writeCsv("a_unique.csv", aUnique, header)
	writeCsv("b_unique.csv", bUnique, header)
	writeCsv("common.csv", common, header)

	fmt.Println("CSV files have been written successfully.")
}

// readCsv はCSVファイルを読み込み、データのスライスとヘッダーを返す
func readCsv(filename string) ([][]string, []string, error) {
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

	header := records[0] // ヘッダーを取得
	data := records[1:]  // ヘッダーを除いたデータ

	return data, header, nil
}

// compareLists は2つのデータセットを比較し、ユニークな行と共通行を返す
func compareLists(a, b [][]string) (aUnique, bUnique, common [][]string) {
	mA := make(map[string]bool)
	mB := make(map[string][]string)

	for _, item := range b {
		mB[item[0]] = item // 'inum'が最初の列であると仮定
	}
	for _, item := range a {
		if bItem, found := mB[item[0]]; found {
			common = append(common, bItem)
		} else {
			aUnique = append(aUnique, item)
		}
	}
	for _, item := range b {
		if _, found := mA[item[0]]; !found {
			bUnique = append(bUnique, item)
		}
	}

	return
}

// writeCsv は指定されたファイル名でデータセットをヘッダー付きでCSVファイルに書き込む
func writeCsv(filename string, data [][]string, header []string) {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// ヘッダーを書き込む
	if err := writer.Write(header); err != nil {
		panic("Error writing header to csv:" + err.Error())
	}

	for _, value := range data {
		if err := writer.Write(value); err != nil {
			panic("Error writing record to csv:" + err.Error())
		}
	}
}
