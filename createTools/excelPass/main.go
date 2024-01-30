package main

import (
	"os"
	"path/filepath"

	"log"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/xuri/excelize/v2"
)

func main() {
	a := app.New()
	w := a.NewWindow("パスワード入力")

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("パスワードを入力")

	submitButton := widget.NewButton("OK", func() {
		password := passwordEntry.Text
		w.Close()
		// パスワードを使用してExcelファイルを保護する
		protectExcelFiles(password)
	})

	w.SetContent(container.NewVBox(
		passwordEntry,
		submitButton,
	))

	w.ShowAndRun()
}

func protectExcelFiles(password string) {
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".xlsx" {
			err := protectSheet(path, password)
			if err != nil {
				log.Printf("ファイル '%s' の保護中にエラーが発生しました: %v\n", path, err)
			} else {
				log.Printf("ファイル '%s' は保護されました\n", path)
			}
		}
		return nil
	})

	if err != nil {
		log.Printf("エラー: %v\n", err)
	}
}

func protectSheet(filePath, password string) error {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return err
	}

	for _, name := range f.GetSheetList() {
		// シートをパスワードで保護
		err := f.ProtectSheet(name, &excelize.SheetProtectionOptions{
			Password: password,
		})
		if err != nil {
			return err
		}
	}

	return f.SaveAs(filePath)
}
