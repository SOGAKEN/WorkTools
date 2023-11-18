package handlers

import (
	"database/sql"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
)

func ConnectDB() *sql.DB {
	// 環境変数や設定ファイルから接続文字列を取得
	connString := getConnectionString()
	db, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatal("データベース接続エラー: ", err)
	}
	return db
}

// getConnectionString - 環境変数や設定ファイルから接続文字列を取得する関数
func getConnectionString() string {
	// TODO: 環境変数や設定ファイルから接続文字列を取得する実装をここに追加
	return "server=サーバー;user id=ユーザー名;password=パスワード;database=データベース名"
}
