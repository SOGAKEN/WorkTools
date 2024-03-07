package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	url := "https://www.yahoo.co.jp/"
	interval := 1 * time.Second // 定期的にリクエストを送る間隔

	for {
		logFile := getLogFileName() // 現在の日付を含むログファイル名を取得

		// ログファイルを開く
		file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}

		statusCode, err := checkURL(url)
		if err != nil {
			// エラーが発生した場合、エラーメッセージをログとコンソールに出力
			logMsg := fmt.Sprintf("[%s]: ERROR %v\n", getCurrentTimeFormatted(), err)
			fmt.Print(logMsg)
			if _, err := file.WriteString(logMsg); err != nil {
				log.Fatalf("Failed to write to log file: %v", err)
			}
		} else if statusCode != 200 {
			// ステータスコードが200以外の場合、ステータスコードを含めてログとコンソールに出力
			logMsg := fmt.Sprintf("[%s]: ERROR status code %d\n", getCurrentTimeFormatted(), statusCode)
			fmt.Print(logMsg)
			if _, err := file.WriteString(logMsg); err != nil {
				log.Fatalf("Failed to write to log file: %v", err)
			}
		} else {
			// ステータスコードが200の場合、成功メッセージをコンソールに出力
			fmt.Printf("[%s]: SUCCESS status code %d\n", getCurrentTimeFormatted(), statusCode)
		}

		file.Close()         // ログファイルを閉じる
		time.Sleep(interval) // 次のリクエストまで待機
	}
}

func getLogFileName() string {
	// 現在の日付をyyyymmdd形式で取得
	currentDate := time.Now().Format("20060102")
	// ログファイル名をフォーマット
	logFile := fmt.Sprintf("./log_%s.csv", currentDate)
	return logFile
}

func getCurrentTimeFormatted() string {
	// 日時を[yyyy-mm-dd hh:mm:ss]形式でフォーマット
	return time.Now().Format("2006-01-02 15:04:05")
}

func checkURL(url string) (int, error) {
	client := http.Client{
		Timeout: 5 * time.Second, // タイムアウト設定
	}
	resp, err := client.Get(url) // HTTP GETリクエストを送信
	if err != nil {
		return 0, err // エラーが発生した場合
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil // ステータスコードを返す
}
