package handlers

import (
	"encoding/csv"
	"io"
	"log"
	"os"
)

func loadProtectedValue() map[string]struct{} {
	file, err := os.Open("protect.csv")
	if err != nil {
		log.Fatal("protect.csvファイルオープンエラー：", err)
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
			log.Fatal("protect.csv読みいこみエラー", err)
		}
		protectedValues[record[0]] = struct{}{}
	}
	return protectedValues
}

func InitCSVWriters() (*csv.Writer, *csv.Writer, *csv.Writer, map[string]struct{}) {
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

	allWriter := csv.NewWriter(allFile)
	listWriter := csv.NewWriter(listFile)
	oneWriter := csv.NewWriter(oneFile)

	addHeaders(allWriter, []string{"logId", "Min-Date", "Max-Date", "How"})
	addHeaders(listWriter, []string{"logId", "Min-Date", "Max-Date", "How"})
	addHeaders(oneWriter, []string{"logId"})

	protectedValues := loadProtectedValue()

	return allWriter, listWriter, oneWriter, protectedValues
}

func addHeaders(writer *csv.Writer, headers []string) {
	err := writer.Write(headers)
	if err != nil {
		log.Fatal("ヘッダーの書き込みエラー", err)
	}
}

func CloseCSVWriters(allWriter, listWriter, oneWriter *csv.Writer) {
	allWriter.Flush()
	listWriter.Flush()
	oneWriter.Flush()
}
