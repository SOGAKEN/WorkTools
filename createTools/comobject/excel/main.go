package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	// 現在のディレクトリを取得
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// ディレクトリ内のファイルを反復処理
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Excelファイルのみを処理
		if !info.IsDir() && (strings.HasSuffix(info.Name(), ".xlsx") || strings.HasSuffix(info.Name(), ".xls")) {
			fmt.Println("Protecting Excel file:", path)

			// PowerShellスクリプトのパス
			psScriptPath := "./protect_excel.ps1"

			// PowerShellスクリプトを実行
			cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-File", psScriptPath, "-excelFilePath", path)
			if err := cmd.Run(); err != nil {
				fmt.Println("Error running PowerShell script:", err)
			} else {
				fmt.Println("Successfully protected:", path)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking through directory: %v\n", err)
	}
}
