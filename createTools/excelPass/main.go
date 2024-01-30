package main

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/xuri/excelize/v2"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("パスワードを引数として指定してください")
        os.Exit(1)
    }
    password := os.Args[1]

    err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if filepath.Ext(path) == ".xlsx" {
            err := setPassword(path, password)
            if err != nil {
                fmt.Printf("ファイル '%s' にパスワードを設定中にエラーが発生しました: %v\n", path, err)
            } else {
                fmt.Printf("ファイル '%s' にパスワードが設定されました\n", path)
            }
        }
        return nil
    })

    if err != nil {
        fmt.Printf("エラー: %v\n", err)
        os.Exit(1)
    }
}

func setPassword(filePath, password string) error {
    f, err := excelize.OpenFile(filePath)
    if err != nil {
        return err
    }
    // excelizeには現在のバージョンで直接的なパスワード設定の機能はないため、
    // ファイルの保護機能など他の方法を検討する必要があります。
    // 以下はファイル保護の例です。
    options := excelize.SetSheetProtectionOptions{
        Password: password,
    }
    for _, name := range f.GetSheetList() {
        f.ProtectSheet(name, &options)
    }
    return f.Save()
}
