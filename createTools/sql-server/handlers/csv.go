package handlers

import (
	"encoding/csv"
	"io"
	"log"
	"os"
)

// protect.csvから保護された値のセットを作成します。
func LoadProtectedValues() map[string]struct{} {
	file, err := os.Open("protect.csv")
	if err != nil {
		log.Fatal("protect.csvファイルオープンエラー: ", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	protectedValues := make(map[string]struct{})

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("protect.csv読み込みエラー: ", err)
		}
		protectedValues[record[0]] = struct{}{}
	}

	return protectedValues
}

// CSVファイルとCSVライターを初期化します。
func InitCSVWriters() (*csv.Writer, *csv.Writer, *csv.Writer) {
	allFile, err := os.Create("all.csv")
	if err != nil {
		log.Fatal("all.csvファイル作成エラー: ", err)
	}

	listFile, err := os.Create("list.csv")
	if err != nil {
		log.Fatal("list.csvファイル作成エラー: ", err)
	}

	oneFile, err := os.Create("one.csv")
	if err != nil {
		log.Fatal("one.csvファイル作成エラー: ", err)
	}

	// CSVライターを作成
	allWriter := csv.NewWriter(allFile)
	listWriter := csv.NewWriter(listFile)
	oneWriter := csv.NewWriter(oneFile)

	// ヘッダーを追加
	addHeaders(allWriter, []string{"ヘッダー1", "ヘッダー2", "ヘッダー3"})
	addHeaders(listWriter, []string{"ヘッダーA", "ヘッダーB", "ヘッダーC"})
	addHeaders(oneWriter, []string{"ヘッダーX", "ヘッダーY", "ヘッダーZ"})

	return allWriter, listWriter, oneWriter
}

// CSVライターにヘッダーを追加します。
func addHeaders(writer *csv.Writer, headers []string) {
	err := writer.Write(headers)
	if err != nil {
		log.Fatal("ヘッダーの書き込みエラー: ", err)
	}
}

// CSVライターをフラッシュして閉じます。
func CloseCSVWriters(allWriter, listWriter, oneWriter *csv.Writer) {
	allWriter.Flush()
	listWriter.Flush()
	oneWriter.Flush()
}
