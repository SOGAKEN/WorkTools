
package main

import (
	"encoding/csv"
	"os"
)

func main() {
	// List Sheet.csv ファイルを読み込む
	listSheetFile, err := os.Open("List Sheet.csv")
	if err != nil {
		panic(err)
	}
	defer listSheetFile.Close()

	// CSV リーダーを作成
	listSheetReader := csv.NewReader(listSheetFile)

	// CSV データを全て読み込む
	listSheetData, err := listSheetReader.ReadAll()
	if err != nil {
		panic(err)
	}

	// Protect.csv ファイルを読み込む
	protectFile, err := os.Open("Protect.csv")
	if err != nil {
		panic(err)
	}
	defer protectFile.Close()

	// CSV リーダーを作成
	protectReader := csv.NewReader(protectFile)

	// CSV データを全て読み込む
	protectData, err := protectReader.ReadAll()
	if err != nil {
		panic(err)
	}

	// Protect.csv の値をマップに保存
	protectMap := make(map[string]bool)
	for _, row := range protectData {
		protectMap[row[0]] = true
	}

	// List Sheet.csv から Protect.csv の値を除外
	var filteredData [][]string
	for _, row := range listSheetData {
		if !protectMap[row[0]] {
			filteredData = append(filteredData, row)
		}
	}

	// 結果を新しいCSVファイルに書き込む
	newFile, err := os.Create("List Sheet_filtered.csv")
	if err != nil {
		panic(err)
	}
	defer newFile.Close()

	writer := csv.NewWriter(newFile)
	err = writer.WriteAll(filteredData)
	if err != nil {
		panic(err)
	}
}
