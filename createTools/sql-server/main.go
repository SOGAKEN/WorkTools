package main

import (
	"log"
	"main/handlers"
	"time"
)

func main() {
	start := time.Now()
	log.Println("プロセス開始")

	db := handlers.ConnectDB()

	// CSVファイルハンドラの初期化
	allWriter, listWriter, oneWriter, protectedValues := handlers.InitCSVWriters()
	defer handlers.CloseCSVWriters(allWriter, listWriter, oneWriter)

	// 最初のクエリの実行
	handlers.ExecuteFirstQuery(db, allWriter, listWriter, protectedValues)

	// 二番目のクエリの実行
	handlers.ExecuteSecondQuery(db, oneWriter, protectedValues)

	elapsed := time.Since(start)
	log.Printf("プロセス終了。所要時間: %s", elapsed)
}
