package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// executePowerShellScript executes the specified PowerShell script with given arguments.
func executePowerShellScript(scriptPath string, args ...string) (string, error) {
	cmdArgs := []string{"-ExecutionPolicy", "Bypass", "-File", scriptPath}
	cmdArgs = append(cmdArgs, args...)

	cmd := exec.Command("powershell", cmdArgs...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func main() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		return
	}

	psScriptPath := filepath.Join(dir, "Set-WordFilesEditProtected.ps1")
	password := "yourPassword"
	editingRestriction := "ReadOnly"
	createNewVersion := "$true" // PowerShell expects $true or $false as string

	// Construct PowerShell script arguments
	args := []string{"-password", password, "-editingRestriction", editingRestriction, "-createNewVersion"}

	// Execute the PowerShell script
	output, err := executePowerShellScript(psScriptPath, args...)
	if err != nil {
		fmt.Printf("Error executing PowerShell script: %v\n", err)
		return
	}

	fmt.Println(output)
}
