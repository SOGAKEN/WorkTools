package handlers

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
)

func loadProtectedValue() (map[string]struct{}, int) {
	file, err := os.Open("protect.csv")
	if err != nil {
		log.Fatal("protect.csvファイルオープンエラー：", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	protectedValues := make(map[string]struct{})
	var secondLineValue int

	lineCount := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("protect.csv読みいこみエラー", err)
		}
		lineCount++
		if lineCount == 2 {
			// Assuming the value is in the first column of the second line
			secondLineValue, err = strconv.Atoi(record[0])
			if err != nil {
				log.Fatal("protect.csv 2行目の値の変換エラー", err)
			}
		}

		if lineCount > 3 {
			protectedValues[record[0]] = struct{}{}
		}
	}
	return protectedValues, secondLineValue
}

func InitCSVWriters() (*csv.Writer, *csv.Writer, *csv.Writer, *csv.Writer, map[string]struct{}, int) {
	allFile, err := os.Create("all.csv")
	if err != nil {
		log.Fatal("all.csvファイル作成エラー: ", err)
	}
	// BOMを追加
	_, err = allFile.WriteString("\uFEFF")
	if err != nil {
		log.Fatal("BOMの書き込みエラー: ", err)
	}

	listFile, err := os.Create("list.csv")
	if err != nil {
		log.Fatal("list.csvファイル作成エラー: ", err)
	}
	// BOMを追加
	_, err = listFile.WriteString("\uFEFF")
	if err != nil {
		log.Fatal("BOMの書き込みエラー: ", err)
	}

	oneFile, err := os.Create("one.csv")
	if err != nil {
		log.Fatal("one.csvファイル作成エラー: ", err)
	}
	// BOMを追加
	_, err = oneFile.WriteString("\uFEFF")
	if err != nil {
		log.Fatal("BOMの書き込みエラー: ", err)
	}

	notProtectFile, err := os.Create("not_protect.csv")
	if err != nil {
		log.Fatal("not_protect.csvファイル作成エラー: ", err)
	}
	_, err = notProtectFile.WriteString("\uFEFF") // Writing BOM
	if err != nil {
		log.Fatal("BOMの書き込みエラー: ", err)
	}

	allWriter := csv.NewWriter(allFile)
	listWriter := csv.NewWriter(listFile)
	oneWriter := csv.NewWriter(oneFile)
	notProtectWriter := csv.NewWriter(notProtectFile)

	addHeaders(allWriter, []string{"logId", "Min-Date", "Max-Date", "How"})
	addHeaders(listWriter, []string{"logId", "Min-Date", "Max-Date", "How", "NOT"})
	addHeaders(oneWriter, []string{"logId", "NOT"})
	addHeaders(notProtectWriter, []string{"LogID"})

	protectedValues, secondLineValue := loadProtectedValue()

	return allWriter, listWriter, oneWriter, notProtectWriter, protectedValues, secondLineValue
}

func addHeaders(writer *csv.Writer, headers []string) {
	err := writer.Write(headers)
	if err != nil {
		log.Fatal("ヘッダーの書き込みエラー", err)
	}
}

func CloseCSVWriters(allWriter, listWriter, oneWriter, notProtectWriter *csv.Writer) {
	allWriter.Flush()
	listWriter.Flush()
	oneWriter.Flush()
	notProtectWriter.Flush()
}
