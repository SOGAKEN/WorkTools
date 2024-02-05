package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// executePowerShellScript は、指定されたPowerShellスクリプトを実行します。
// スクリプトへの引数もサポートされます。
func executePowerShellScript(scriptPath string, args ...string) (string, error) {
	// PowerShellコマンドの実行ポリシーをBypassに設定し、セキュリティ警告を回避
	cmdArgs := []string{"-ExecutionPolicy", "Bypass", "-File", scriptPath}
	cmdArgs = append(cmdArgs, args...)

	cmd := exec.Command("powershell", cmdArgs...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func main() {
	// 現在のディレクトリを取得
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		return
	}

	psScriptPath := filepath.Join(dir, "Set-WordFilesEditProtected.ps1")
	password := "yourPassword"

	// PowerShellスクリプトの引数
	args := []string{"-password", password}

	// PowerShellスクリプトを実行
	output, err := executePowerShellScript(psScriptPath, args...)
	if err != nil {
		fmt.Printf("Error executing PowerShell script: %v\n", err)
		return
	}

	fmt.Println(output)
}
