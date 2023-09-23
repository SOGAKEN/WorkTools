package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/xuri/excelize/v2"
)

func printSheetData(f *excelize.File, sheetName string) {
	rows, _ := f.GetRows(sheetName)
	for _, row := range rows {
		for _, cell := range row {
			fmt.Print(cell, "\t")
		}
		fmt.Println()
	}
}

func colNumToName(col int) string {
	name := ""
	for col > 0 {
		col--
		name = string('A'+col%26) + name
		col /= 26
	}
	return name
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("使用方法: go run main.go <excelファイル>")
		os.Exit(1)
	}

	filePath := os.Args[1]

	f, err := excelize.OpenFile(filePath)
	if err != nil {
		panic(err)
	}

	sheetName := f.GetSheetName(0)

	// 1. エクセルファイルの内容をログに出力
	fmt.Println("=== Original Excel Data ===")
	printSheetData(f, sheetName)

	maxCols, _ := f.GetCols(sheetName)
	newColName := colNumToName(len(maxCols) + 1)

	// 新しいヘッダーを追加
	f.SetCellValue(sheetName, newColName+"1", "新しいヘッダー")

	rows, _ := f.GetRows(sheetName)
	for rowIndex, row := range rows {
		if rowIndex == 0 {
			continue // ヘッダー行をスキップ
		}
		val, err := strconv.Atoi(row[1])
		if err == nil {
			// 値が整数の場合
			doubleVal := val * 2
			cellName := fmt.Sprintf("%s%d", newColName, rowIndex+1)
			f.SetCellValue(sheetName, cellName, doubleVal)
		} else {
			// 値が整数でない場合
			cellName := fmt.Sprintf("%s%d", newColName, rowIndex+1)
			f.SetCellValue(sheetName, cellName, row[1])
		}
	}

	// 2. 変更後のデータをログに出力
	fmt.Println("\n=== Updated Excel Data ===")
	printSheetData(f, sheetName)

	if err := f.SaveAs("output.xlsx"); err != nil {
		panic(err)
	}
}
