package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/xuri/excelize/v2"
)

func init() {
	// ログの出力設定を初期化
	log.SetOutput(os.Stderr)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// checkSheetProtection はZIPファイル内の特定のファイル（Excelシート）でsheetProtectionタグを探します。
// sheetProtectionが見つかればtrue、見つからなければfalseを返します。
func checkSheetProtection(filePath string) bool {
	zipReader, err := zip.OpenReader(filePath)
	if err != nil {
		log.Printf("Error opening file: %v\n", err)
		return false
	}
	defer zipReader.Close()

	for _, zipFile := range zipReader.File {
		if strings.Contains(zipFile.Name, "xl/worksheets/") {
			f, err := zipFile.Open()
			if err != nil {
				log.Printf("Error opening zip file: %v\n", err)
				continue
			}

			content, err := io.ReadAll(f)
			f.Close()
			if err != nil {
				log.Printf("Error reading file content: %v\n", err)
				continue
			}

			if strings.Contains(string(content), "<sheetProtection") {
				return true
			}
		}
	}
	return false
}

func protectExcelFiles(filePath, password string) {
	// Excelファイルを開く
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		log.Printf("ファイル %s を開く際にエラーが発生しました: %v\n", filePath, err)
		return
	}

	for _, name := range f.GetSheetMap() {
		// ブックの保護を設定
		if err := f.ProtectSheet(name, &excelize.SheetProtectionOptions{
			Password:            password,
			SelectLockedCells:   true,
			SelectUnlockedCells: true,
		}); err != nil {
			log.Printf("ワークブック %s の保護設定時にエラーが発生しました: %v\n", filePath, err)
			return
		}
	}

	// 保護を適用したファイルを保存
	if err := f.Save(); err != nil {
		log.Printf("保護されたワークブック %s を保存する際にエラーが発生しました: %v\n", filePath, err)
		return
	}

	log.Printf("保護されました: %s\n", filePath)
}

func main() {
	dirPath := "./"             // 対象のディレクトリ
	password := "your_password" // 保護に使用するパスワード

	filepath.Walk(dirPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("パス %q の走査中にエラーが発生しました: %v\n", filePath, err)
			return err
		}
		if !info.IsDir() && filepath.Ext(filePath) == ".xlsx" {
			if !checkSheetProtection(filePath) {
				// sheetProtectionが見つからなければ、保護処理を行う
				protectExcelFiles(filePath, password)
			} else {
				log.Printf("保護の必要なし: %s\n", filePath)
			}
		}
		return nil
	})
}
