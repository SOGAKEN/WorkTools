package main

import (
	"log"
	"main/handlers"
	"main/model"
	"time"
)

func main() {
	start := time.Now()
	log.Println("プロセス開始")

	db := handlers.ConnectDB()

	// CSVファイルハンドラの初期化
	allWriter, listWriter, oneWriter, notProtectWriter, protectedValues, secondLineValue := handlers.InitCSVWriters()
	defer handlers.CloseCSVWriters(allWriter, listWriter, oneWriter, notProtectWriter)

	// 最初のクエリの実行
	params := model.FirstQueryParams{
		DB:               db,
		AllWriter:        allWriter,
		ListWriter:       listWriter,
		NotProtectWriter: notProtectWriter,
		ProtectedValues:  protectedValues,
		SecondLineValue:  secondLineValue,
	}

	handlers.ExecuteFirstQuery(params)

	// 二番目のクエリの実行
	params2 := model.SecondQueryParams{
		DB:              db,
		OneWriter:       oneWriter,
		ProtectedValues: protectedValues,
	}
	handlers.ExecuteSecondQuery(params2)

	elapsed := time.Since(start)
	log.Printf("プロセス終了。所要時間: %s", elapsed)
}
