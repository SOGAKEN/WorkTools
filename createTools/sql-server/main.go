package main

import (
	"database/sql"
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
)

type DataRecord struct {
	LogID                string
	MinDate              string
	MaxDate              string
	HowManyDaysFromToday int
}

func main() {
	start := time.Now()
	log.Println("プロセス開始")

	// データベースに接続
	connString := "server=あなたのサーバー;user id=あなたのユーザー名;password=あなたのパスワード;database=あなたのデータベース"
	db, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatal("接続プールの作成エラー: ", err.Error())
	}
	defer db.Close()

	// 最初のクエリの実行
	executeFirstQuery(db)

	// 二番目のクエリの実行
	executeSecondQuery(db)

	elapsed := time.Since(start)
	log.Printf("プロセス終了。所要時間: %s", elapsed)
}

func executeFirstQuery(db *sql.DB) {
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
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal("クエリ実行エラー: ", err.Error())
	}
	defer rows.Close()

	allFile, err := os.Create("all.csv")
	if err != nil {
		log.Fatal("all.csvファイル作成不可:", err)
	}
	defer allFile.Close()
	allWriter := csv.NewWriter(allFile)
	defer allWriter.Flush()

	listFile, err := os.Create("list.csv")
	if err != nil {
		log.Fatal("list.csvファイル作成不可:", err)
	}
	defer listFile.Close()
	listWriter := csv.NewWriter(listFile)
	defer listWriter.Flush()

	for rows.Next() {
		var record DataRecord
		err := rows.Scan(&record.LogID, &record.MinDate, &record.MaxDate, &record.HowManyDaysFromToday)
		if err != nil {
			log.Fatal("行のスキャンエラー: ", err.Error())
		}

		allWriter.Write([]string{record.LogID, record.MinDate, record.MaxDate, strconv.Itoa(record.HowManyDaysFromToday)})

		if record.HowManyDaysFromToday >= 183 {
			listWriter.Write([]string{record.LogID, record.MinDate, record.MaxDate, strconv.Itoa(record.HowManyDaysFromToday)})
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal("行の読み取り中のエラー: ", err)
	}
}

func executeSecondQuery(db *sql.DB) {
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

	oneFile, err := os.Create("one.csv")
	if err != nil {
		log.Fatal("one.csvファイル作成不可:", err)
	}
	defer oneFile.Close()
	oneWriter := csv.NewWriter(oneFile)
	defer oneWriter.Flush()

	for rows.Next() {
		var value string
		err := rows.Scan(&value)
		if err != nil {
			log.Fatal("行のスキャンエラー: ", err.Error())
		}

		oneWriter.Write([]string{value})
	}

	if err := rows.Err(); err != nil {
		log.Fatal("行の読み取り中のエラー: ", err)
	}
}
