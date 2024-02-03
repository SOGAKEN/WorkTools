package main

import (
	"fmt"
	"io/ioutil"
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

	// ディレクトリ内のファイルリストを取得
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	// PowerShellスクリプトのパス
	psScriptPath := "./protect_excel.ps1"

	// Excelファイルのみをフィルタリングしてブック保護をかける
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".xlsx") || strings.HasSuffix(file.Name(), ".xls") {
			excelFilePath := filepath.Join(dir, file.Name())
			fmt.Println("Protecting Excel file:", excelFilePath)

			// PowerShellスクリプトを実行
			cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-File", psScriptPath, "-excelFilePath", excelFilePath)
			if err := cmd.Run(); err != nil {
				fmt.Println("Error running PowerShell script:", err)
			} else {
				fmt.Println("Successfully protected:", excelFilePath)
			}
		}
	}
}
