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
    // 実行ファイルの絶対パスを取得
    exePath, err := os.Executable()
    if err != nil {
        fmt.Printf("Error getting executable path: %v\n", err)
        return
    }
    dir := filepath.Dir(exePath)

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
