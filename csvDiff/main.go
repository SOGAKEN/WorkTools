package main

import (
	"encoding/csv"
	"os"
	"path/filepath"
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
			processCSV(entry.Name())
		}
	}
}

func processCSV(fileName string) {
	// CSVファイルを開く
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// CSVリーダーを作成
	reader := csv.NewReader(file)

	// CSVデータを読み込む
	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	// ユニークと重複データを保存するマップ
	uniqueRecords := make(map[string][]string)
	duplicateRecords := make([][]string, 0)

	// レコードをチェック
	for _, record := range records {
		if _, exists := uniqueRecords[record[0]]; exists {
			// 重複レコードに追加
			duplicateRecords = append(duplicateRecords, record)
			delete(uniqueRecords, record[0])
		} else if !contains(duplicateRecords, record) {
			// ユニークレコードに追加
			uniqueRecords[record[0]] = record
		}
	}

	// ユニークなレコードのみのCSVを作成
	writeCSV(fileName+"_unique.csv", uniqueRecords)

	// 重複するレコードのCSVを作成
	writeCSV(fileName+"_only.csv", duplicateRecords)
}

func writeCSV(fileName string, records interface{}) {
	// 新しいCSVファイルを開く
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// CSVライターを作成
	writer := csv.NewWriter(file)
	defer writer.Flush()

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
