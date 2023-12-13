package handlers

import (
	"database/sql"
	"encoding/csv"
	"log"
	"main/model"
	"strconv"
)

func ExecuteFirstQuery(db *sql.DB, allWriter, listWriter *csv.Writer, protectedValues map[string]struct{}) {
	query := `...` // 以前と同じクエリ
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal("クエリ実行エラー: ", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var record model.DataRecord
		err := rows.Scan(&record.LogID, &record.MinDate, &record.MaxDate, &record.HowManyDaysFromToday)
		if err != nil {
			log.Fatal("行のスキャンエラー: ", err.Error())
		}

		// 保護された値をチェックして除外
		if _, ok := protectedValues[record.LogID]; ok {
			continue
		}

		// CSVファイルに書き込み
		allWriter.Write([]string{record.LogID, record.MinDate, record.MaxDate, strconv.Itoa(record.HowManyDaysFromToday)})
		if record.HowManyDaysFromToday >= 183 {
			listWriter.Write([]string{record.LogID, record.MinDate, record.MaxDate, strconv.Itoa(record.HowManyDaysFromToday)})
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal("行の読み取り中のエラー: ", err)
	}
}

func ExecuteSecondQuery(db *sql.DB, oneWriter *csv.Writer, protectedValues map[string]struct{}) {
	query := `...` // 以前と同じクエリ
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal("クエリ実行エラー: ", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var value string
		err := rows.Scan(&value)
		if err != nil {
			log.Fatal("行のスキャンエラー: ", err.Error())
		}

		// 保護された値をチェックして除外
		if _, ok := protectedValues[value]; ok {
			continue
		}

		// CSVファイルに書き込み
		oneWriter.Write([]string{value})
	}

	if err := rows.Err(); err != nil {
		log.Fatal("行の読み取り中のエラー: ", err)
	}
}
