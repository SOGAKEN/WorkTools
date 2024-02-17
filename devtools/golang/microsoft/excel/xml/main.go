package main

import (
    "archive/zip"
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "strings"
    "time"

    "github.com/xuri/excelize/v2"
)

func init() {
    log.SetOutput(os.Stderr)
    log.SetFlags(0) // ログの日付や時間の出力を無効化
}

func checkSheetProtection(filePath string) (bool, error) {
    zipReader, err := zip.OpenReader(filePath)
    if err != nil {
        return false, err
    }
    defer zipReader.Close()

    for _, zipFile := range zipReader.File {
        if strings.Contains(zipFile.Name, "xl/worksheets/") {
            f, err := zipFile.Open()
            if err != nil {
                return false, err
            }

            content, err := io.ReadAll(f)
            f.Close()
            if err != nil {
                return false, err
            }

            if strings.Contains(string(content), "<sheetProtection") {
                return true, nil
            }
        }
    }
    return false, nil
}

func protectExcelFiles(filePath, password string) error {
    f, err := excelize.OpenFile(filePath)
    if err != nil {
        return err
    }

    for _, name := range f.GetSheetMap() {
        err := f.ProtectSheet(name, &excelize.SheetProtectionOptions{
            Password:            password,
            SelectLockedCells:   true,
            SelectUnlockedCells: true,
        })
        if err != nil {
            return err
        }
    }

    return f.Save()
}

func logProgress(fileCount int, result, filePath string) {
    timestamp := time.Now().Format("15:04:05")
    fmt.Printf("[%s][%d][%s] | %s\n", timestamp, fileCount, result, filePath)
}

func main() {
    dirPath := "./"
    password := "your_password"
    fileCount := 0

    filepath.Walk(dirPath, func(filePath string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.IsDir() && filepath.Ext(filePath) == ".xlsx" {
            fileCount++
            protected, err := checkSheetProtection(filePath)

            if err != nil {
                logProgress(fileCount, "NG", filePath)
                return nil // エラーを記録して次のファイルに進む
            }

            if protected {
                logProgress(fileCount, "PASS", filePath)
            } else {
                err := protectExcelFiles(filePath, password)
                if err != nil {
                    logProgress(fileCount, "NG", filePath)
                } else {
                    logProgress(fileCount, "OK", filePath)
                }
            }
        }
        return nil
    })

    fmt.Println("処理が完了しました。エンターキーを押して終了してください...")
    fmt.Scanln()
}
