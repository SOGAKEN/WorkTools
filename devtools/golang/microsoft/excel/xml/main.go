package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func checkSheetProtection(filePath string) {
	fmt.Printf("Checking: %s\n", filePath)

	zipReader, err := zip.OpenReader(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer zipReader.Close()

	for _, zipFile := range zipReader.File {
		if strings.Contains(zipFile.Name, "xl/worksheets/") {
			f, err := zipFile.Open()
			if err != nil {
				fmt.Println("Error opening zip file:", err)
				continue
			}

			content, err := io.ReadAll(f)
			f.Close()
			if err != nil {
				fmt.Println("Error reading file content:", err)
				continue
			}

			if strings.Contains(string(content), "<sheetProtection") {
				fmt.Printf("SheetProtection tag found in %s\n", zipFile.Name)
			}
		}
	}
}

func walkAndCheckFiles(rootPath string) {
	filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path %q: %v\n", path, err)
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".xlsx") {
			checkSheetProtection(path)
		}
		return nil
	})
}

func main() {
	walkAndCheckFiles(".")
}
