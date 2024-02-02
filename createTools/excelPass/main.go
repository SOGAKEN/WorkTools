package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/xuri/excelize/v2"
)

func init() {
	// ログの出力設定を初期化
	log.SetOutput(os.Stderr)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func protectExcelFiles(path, password string) error {
	excelFilesFound := false // Excelファイルの存在フラグ

	// ディレクトリ内の全ての.xlsxファイルを検索
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("パス %q の走査中にエラーが発生しました: %v\n", filePath, err)
			return err
		}
		if !info.IsDir() && filepath.Ext(filePath) == ".xlsx" {
			excelFilesFound = true // Excelファイルを見つけたらフラグをtrueに設定

			// Excelファイルを開く
			f, err := excelize.OpenFile(filePath)
			if err != nil {
				log.Printf("ファイル %s を開く際にエラーが発生しました: %v\n", filePath, err)
				return err
			}

			// ブックの保護を設定
			if err := f.ProtectWorkbook(&excelize.WorkbookProtectionOptions{
				Password:      password,
				LockStructure: true,
				LockWindows:   true,
			}); err != nil {
				log.Printf("ワークブック %s の保護設定時にエラーが発生しました: %v\n", filePath, err)
				return err
			}
			// 保護を適用したファイルを保存
			if err := f.Save(); err != nil {
				log.Printf("保護されたワークブック %s を保存する際にエラーが発生しました: %v\n", filePath, err)
				return err
			}
			log.Printf("保護されました: %s\n", filePath)
		}
		return nil
	})

	if err != nil {
		return err
	}

	// Excelファイルが一つも見つからなかった場合のログ出力
	if !excelFilesFound {
		log.Println("エクセルファイルがありません。")
	}

	return nil
}

func main() {
	// 対象のディレクトリとパスワードを設定
	dirPath := "./"             // ここには実行ファイルと同じ階層のディレクトリを指定
	password := "your_password" // 保護に使用するパスワードを指定

	if err := protectExcelFiles(dirPath, password); err != nil {
		log.Printf("Excelファイルの保護中にエラーが発生しました: %v\n", err)
	}
}
