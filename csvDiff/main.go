package main

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// 指定フォルダのファイル一覧を取得
	entries, err := os.ReadDir(".")
	if err != nil {
		panic(err)
	}

	// CSVファイルを処理
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".csv" {
			processCSV(entry.Name(), "switch_call_id")
		}
	}
}

func processCSV(fileName string, columnName string) {
	// CSVファイルを開く
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// CSVリーダーを作成
	reader := csv.NewReader(file)

	// 最初の行を読み込む（ヘッダーとして）
	header, err := reader.Read()
	if err != nil {
		panic(err)
	}

	// カラム名からインデックスを見つける
	columnIndex := -1
	for i, col := range header {
		if strings.ToLower(col) == strings.ToLower(columnName) {
			columnIndex = i
			break
		}
	}
	if columnIndex == -1 {
		panic("Column " + columnName + " not found in file " + fileName)
	}

	// 残りのCSVデータを読み込む
	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	// レコードのカウントを保存するマップ
	counts := make(map[string]int)
	for _, record := range records {
		if len(record) > columnIndex {
			counts[record[columnIndex]]++
		}
	}

	// ユニークと重複レコードを分類
	uniqueRecords := make([][]string, 0)
	duplicateRecords := make([][]string, 0)
	for _, record := range records {
		if len(record) > columnIndex {
			if counts[record[columnIndex]] == 1 {
				uniqueRecords = append(uniqueRecords, record)
			} else {
				duplicateRecords = append(duplicateRecords, record)
			}
		}
	}

	// ユニークなレコードのみのCSVを作成（ヘッダーを含む）
	writeCSV(fileName+"_unique.csv", header, uniqueRecords)

	// 重複するレコードのCSVを作成（ヘッダーを含む）
	writeCSV(fileName+"_duplicates.csv", header, duplicateRecords)
}

func writeCSV(fileName string, header []string, records [][]string) {
	// 新しいCSVファイルを開く
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// CSVライターを作成
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// ヘッダーを書き込む
	if err := writer.Write(header); err != nil {
		panic(err)
	}

	// レコードを書き込む
	for _, record := range records {
		if err := writer.Write(record); err != nil {
			panic(err)
		}
	}
}
