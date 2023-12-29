package handlers

import (
	"database/sql"
	"encoding/csv"
	"log"
	"main/model"
	"strconv"
	"strings"
)

func ExecuteFirstQuery(params model.FirstQueryParams) {
	query := `
        SELECT 
            logid,
            MIN(row_date) MinDate,
            MAX(row_date) MaxDate,
            DATEDIFF(day, MAX(row_date), GETDATE()) HowManyDaysFromToday 
        FROM 
            dagent 
        GROUP BY 
            logid 
        ORDER BY 
            logid
    `
	rows, err := params.DB.Query(query)
	if err != nil {
		log.Fatal("クエリ実行エラー: ", err.Error())
	}
	defer rows.Close()

	remainingProtected := make(map[string]struct{}, len(params.ProtectedValues))
	for k := range params.ProtectedValues {
		remainingProtected[k] = struct{}{}
	}

	for rows.Next() {
		var record model.DataRecord
		err := rows.Scan(&record.LogID, &record.MinDate, &record.MaxDate, &record.HowManyDaysFromToday)
		if err != nil {
			log.Fatal("行のスキャンエラー: ", err.Error())
		}

		trimedLogID := strings.TrimSpace(record.LogID)

		params.AllWriter.Write([]string{record.LogID, record.MinDate, record.MaxDate, strconv.Itoa(record.HowManyDaysFromToday)})

		if _, ok := params.ProtectedValues[trimedLogID]; ok {
			// protectedValues に含まれる場合、新しいカラムを追加して書き込み
			delete(remainingProtected, trimedLogID)
			params.ListWriter.Write([]string{record.LogID, record.MinDate, record.MaxDate, strconv.Itoa(record.HowManyDaysFromToday), "not_End"})
		} else if record.HowManyDaysFromToday >= params.SecondLineValue {
			// protectedValues に含まれず、かつ 183 日以上前の場合
			params.ListWriter.Write([]string{record.LogID, record.MinDate, record.MaxDate, strconv.Itoa(record.HowManyDaysFromToday)})
		}
	}

	for id := range remainingProtected {
		params.NotProtectWriter.Write([]string{id})
	}

	if err := rows.Err(); err != nil {
		log.Fatal("行の読み取り中のエラー: ", err)
	}
}

func ExecuteSecondQuery(db *sql.DB, oneWriter *csv.Writer, protectedValues map[string]struct{}) {
	query := `
        (SELECT DISTINCT value FROM agent)
        EXCEPT
        (SELECT DISTINCT logid FROM dagent)
    `
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

		trimedValue := strings.TrimSpace(value)
		if _, ok := protectedValues[trimedValue]; ok {
			oneWriter.Write([]string{trimedValue, "NOT"})
		} else {
			oneWriter.Write([]string{trimedValue})
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal("行の読み取り中のエラー: ", err)
	}
}
