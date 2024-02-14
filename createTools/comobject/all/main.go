package main

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
)

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
    exePath, err := os.Executable()
    if err != nil {
        fmt.Printf("Error getting executable path: %v\n", err)
        return
    }
    dir := filepath.Dir(exePath)
    psScriptPath := filepath.Join(dir, "Set-FileEditProtected.ps1")
    outputCsv := filepath.Join(dir, "output.csv")
    password := "YourHardcodedPassword" // パスワードをここにハードコーディング

    fmt.Println("Processing files...")
    output, err := executePowerShellScript(psScriptPath, dir, outputCsv, password)
    if err != nil {
        fmt.Printf("Error executing PowerShell script: %v\n", err)
        fmt.Println("Output:", output)
        return
    }

    fmt.Println("Completed. Check the output CSV for details.")
}
