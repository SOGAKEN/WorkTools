package handlers

import (
	"encoding/csv"
	"log"
	"os"
)

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

	return csv.NewWriter(allFile), csv.NewWriter(listFile), csv.NewWriter(oneFile)
}

func CloseCSVWriters(allWriter, listWriter, oneWriter *csv.Writer) {
	allWriter.Flush()
	listWriter.Flush()
	oneWriter.Flush()
}
