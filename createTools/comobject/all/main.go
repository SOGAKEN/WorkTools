package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// executePowerShellScript executes the specified PowerShell script with given arguments.
func executePowerShellScript(scriptPath, directoryPath, outputCsv, password string) (string, error) {
	cmdArgs := []string{
		"-ExecutionPolicy", "Bypass",
		"-NoProfile",
		"-File", scriptPath,
		"-directoryPath", directoryPath,
		"-outputCsv", outputCsv,
		"-password", password,
	}

	cmd := exec.Command("powershell", cmdArgs...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func main() {
	// Go実行ファイルと同じディレクトリを取得
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		return
	}

	// PowerShellスクリプトのパス
	psScriptPath := filepath.Join(dir, "Set-FileEditProtected.ps1")

	// 出力CSVファイルのパス
	outputCsv := filepath.Join(dir, "output.csv")

	// パスワード
	password := "YourPassword"

	// PowerShellスクリプトの実行
	fmt.Println("Processing files...")
	output, err := executePowerShellScript(psScriptPath, dir, outputCsv, password)
	if err != nil {
		fmt.Printf("Error executing PowerShell script: %v\n", err)
		fmt.Println("Output:", output)
		return
	}

	fmt.Println("Completed. Check the output CSV for details.")
}
