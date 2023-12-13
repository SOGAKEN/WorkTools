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

	// CSVファイルハンドラと保護された値のセットの初期化
	allWriter, listWriter, oneWriter, protectedValues := handlers.InitCSVWriters()
	defer handlers.CloseCSVWriters(allWriter, listWriter, oneWriter)

	// 最初のクエリの実行
	handlers.ExecuteFirstQuery(db, allWriter, listWriter, protectedValues)

	// 二番目のクエリの実行（ExecuteSecondQuery関数の実装に基づいて修正が必要）
	handlers.ExecuteSecondQuery(db, oneWriter, protectedValues) // 仮の呼び出し

	elapsed := time.Since(start)
	log.Printf("プロセス終了。所要時間: %s", elapsed)
}
