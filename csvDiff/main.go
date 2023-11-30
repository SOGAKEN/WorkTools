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
			processCSV(entry.Name(), "inum")
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

	// ユニークと重複データを保存するマップ
	uniqueRecords := make(map[string][]string)
	duplicateRecords := make([][]string, 0)

	// レコードをチェック
	for _, record := range records {
		if len(record) <= columnIndex {
			continue // カラム数が足りないレコードは無視
		}
		key := record[columnIndex]
		if _, exists := uniqueRecords[key]; exists {
			// 重複レコードに追加
			duplicateRecords = append(duplicateRecords, record)
			delete(uniqueRecords, key)
		} else if !contains(duplicateRecords, record) {
			// ユニークレコードに追加
			uniqueRecords[key] = record
		}
	}

	// ユニークなレコードのみのCSVを作成（ヘッダーを含む）
	writeCSV(fileName+"_unique.csv", header, uniqueRecords)

	// 重複するレコードのCSVを作成（ヘッダーを含む）
	writeCSV(fileName+"_only.csv", header, duplicateRecords)
}

func writeCSV(fileName string, header []string, records interface{}) {
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
	switch v := records.(type) {
	case map[string][]string:
		for _, record := range v {
			if err := writer.Write(record); err != nil {
				panic(err)
			}
		}
	case [][]string:
		for _, record := range v {
			if err := writer.Write(record); err != nil {
				panic(err)
			}
		}
	default:
		panic("Invalid record type")
	}
}

func contains(records [][]string, record []string) bool {
	for _, r := range records {
		if equal(r, record) {
			return true
		}
	}
	return false
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
