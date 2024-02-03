package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	// 現在のディレクトリを取得
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	password := "yourPassword"
	psScriptPath := filepath.Join(dir, "Set-WordFilesEditProtected.ps1")

	// PowerShellスクリプトを実行
	cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-File", psScriptPath, "-password", password)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(output))
}
