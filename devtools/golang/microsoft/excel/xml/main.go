package main

import (
    "archive/zip"
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "strings"

    "github.com/xuri/excelize/v2"
)

func init() {
    log.SetOutput(os.Stderr)
    log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func checkSheetProtection(filePath string) bool {
    zipReader, err := zip.OpenReader(filePath)
    if err != nil {
        log.Printf("Error opening file: %v\n", err)
        return false
    }
    defer zipReader.Close()

    for _, zipFile := range zipReader.File {
        if strings.Contains(zipFile.Name, "xl/worksheets/") {
            f, err := zipFile.Open()
            if err != nil {
                log.Printf("Error opening zip file: %v\n", err)
                continue
            }

            content, err := io.ReadAll(f)
            f.Close()
            if err != nil {
                log.Printf("Error reading file content: %v\n", err)
                continue
            }

            if strings.Contains(string(content), "<sheetProtection") {
                return true
            }
        }
    }
    return false
}

func protectExcelFiles(filePath, password string) {
    f, err := excelize.OpenFile(filePath)
    if err != nil {
        log.Printf("ファイル %s を開く際にエラーが発生しました: %v\n", filePath, err)
        return
    }

    for _, name := range f.GetSheetMap() {
        if err := f.ProtectSheet(name, &excelize.SheetProtectionOptions{
            Password:            password,
            SelectLockedCells:   true,
            SelectUnlockedCells: true,
        }); err != nil {
            log.Printf("ワークブック %s の保護設定時にエラーが発生しました: %v\n", filePath, err)
            return
        }
    }

    if err := f.Save(); err != nil {
        log.Printf("保護されたワークブック %s を保存する際にエラーが発生しました: %v\n", filePath, err)
        return
    }

    log.Printf("保護されました: %s\n", filePath)
}

func main() {
    dirPath := "./"
    password := "your_password"
    fileCount := 0

    filepath.Walk(dirPath, func(filePath string, info os.FileInfo, err error) error {
        if err != nil {
            log.Printf("パス %q の走査中にエラーが発生しました: %v\n", filePath, err)
            return err
        }
        if !info.IsDir() && filepath.Ext(filePath) == ".xlsx" {
            fileCount++
            fmt.Printf("Processing file %d: %s\n", fileCount, filePath)
            if !checkSheetProtection(filePath) {
                protectExcelFiles(filePath, password)
            } else {
                log.Printf("保護の必要なし: %s\n", filePath)
            }
        }
        return nil
    })

    fmt.Println("処理が完了しました。エンターキーを押して終了してください...")
    fmt.Scanln()
}
