package main

import (
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
	fmt.Println(exePath)
	// backupディレクトリのパスを定義
	baseDir := filepath.Dir(exePath)
	backupDir := filepath.Join(baseDir, "backup")
	// backupディレクトリがなければ作成
	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		if err := os.Mkdir(backupDir, 0755); err != nil {
			panic(err)
		}
	}

	// ファイルシステムを再帰的に探索し、ファイルをコピー
	err = filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if os.IsPermission(err) {
				// アクセス許可エラーが発生した場合、エラーをログに記録し、処理を続行する
				fmt.Println("アクセス許可エラー:", path)
				return nil
			}
			return err
		}
		// ドットから始まるファイルやフォルダをスキップ
		if strings.HasPrefix(filepath.Base(path), ".") {
			if info.IsDir() {
				return filepath.SkipDir // ディレクトリの場合、そのディレクトリをスキップ
			}
			return nil // ファイルの場合、単に次へ進む
		}

		// ディレクトリはスキップ（ただし、backupディレクトリ自体は除外）
		if info.IsDir() {
			return nil
		}

		// backupディレクトリ内のファイルはスキップ
		if strings.HasPrefix(path, backupDir) {
			return nil
		}

		// 元のファイルパスからbackupディレクトリへの相対パスを計算
		relPath, err := filepath.Rel(baseDir, path)
		if err != nil {
			return err
		}

		// ファイル名に(bak)を追加
		dir, file := filepath.Split(relPath)
		destDir := filepath.Join(backupDir, dir)
		ext := filepath.Ext(file)
		nameWithoutExt := file[:len(file)-len(ext)]
		destPath := filepath.Join(destDir, nameWithoutExt+"-bakup"+ext)

		// 必要なサブディレクトリを作成
		if err := os.MkdirAll(destDir, 0755); err != nil {
			return err
		}

		// ファイルをコピー
		return copyFile(path, destPath)
	})
	if err != nil {
		panic(err)
	}
}

// copyFile は、srcファイルをdestファイルにコピーします。
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
