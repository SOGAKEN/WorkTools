package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// 実行ファイルのディレクトリを取得
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	baseDir := filepath.Dir(exePath)
	backupDir := filepath.Join(baseDir, "backup")
	// CSVファイルのパスを定義
	csvFilePath := filepath.Join(baseDir, "backup_results.csv")

	// backupディレクトリがなければ作成
	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		if err := os.Mkdir(backupDir, 0755); err != nil {
			panic(err)
		}
	}

	// CSVファイルを開く（なければ作成）
	csvFile, err := os.OpenFile(csvFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	// BOMを書き込む
	if _, err := csvFile.Write([]byte{0xEF, 0xBB, 0xBF}); err != nil {
		panic(err)
	}

	// CSVライターを作成
	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()

	// ファイルシステムを再帰的に探索し、ファイルをコピー
	err = filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if os.IsPermission(err) {
				fmt.Println("アクセス許可エラー:", path)
				return nil
			}
			return err
		}
		if strings.HasPrefix(filepath.Base(path), ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if info.IsDir() || strings.HasPrefix(path, backupDir) {
			return nil
		}

		relPath, err := filepath.Rel(baseDir, path)
		if err != nil {
			return err
		}

		dir, file := filepath.Split(relPath)
		destDir := filepath.Join(backupDir, dir)
		ext := filepath.Ext(file)
		nameWithoutExt := file[:len(file)-len(ext)]
		destPath := filepath.Join(destDir, nameWithoutExt+"-backup"+ext)

		if err := os.MkdirAll(destDir, 0755); err != nil {
			return err
		}

		// ファイルをコピーし、結果をCSVに記録
		err = copyFile(path, destPath)
		result := "OK"
		if err != nil {
			result = "NG"
		}
		if writeErr := csvWriter.Write([]string{file, destPath, result}); writeErr != nil {
			return writeErr
		}

		return err
	})
	if err != nil {
		panic(err)
	}
}

func copyFile(src, dest string) error {
	input, err := os.Open(src)
	if err != nil {
		return err
	}
	defer input.Close()

	output, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer output.Close()

	_, err = io.Copy(output, input)
	return err
}
