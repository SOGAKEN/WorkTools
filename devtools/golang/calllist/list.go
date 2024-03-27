package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/xuri/excelize/v2"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <ExcelFilePath>")
		return
	}

	filePath := os.Args[1]
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	// A列のセルを順に確認し、空のセルに達するまで処理を続ける
	for i := 2; ; i++ { // 1行目はヘッダーなので2から開始
		cellA, _ := f.GetCellValue("Sheet1", fmt.Sprintf("A%d", i))
		if cellA == "" {
			break // A列のセルが空なら処理を終了
		}

		cellB, _ := f.GetCellValue("Sheet1", fmt.Sprintf("B%d", i))
		cellC, _ := f.GetCellValue("Sheet1", fmt.Sprintf("C%d", i))

		switch cellA {
		case "In Range":
			start, _ := strconv.Atoi(cellB)
			end, _ := strconv.Atoi(cellC)
			for j := start; j <= end; j++ {
				// 生成するリストの位置を調整するために、J列に出力する際の行番号を計算
				rowIndex := i + (j - start)
				f.SetCellInt("Sheet1", fmt.Sprintf("J%d", rowIndex), j)
			}
		case "Equal To":
			value, _ := strconv.Atoi(cellB)
			f.SetCellInt("Sheet1", fmt.Sprintf("J%d", i), value)
		}
	}

	// 変更を保存
	if err := f.SaveAs(filePath); err != nil {
		fmt.Println(err)
	}
}
