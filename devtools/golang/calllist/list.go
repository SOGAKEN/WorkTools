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

	var results []int // 結果を格納するためのスライス

	for i := 1; ; i++ { // A列が空になるまで繰り返し
		cellA, _ := f.GetCellValue("Sheet1", fmt.Sprintf("A%d", i))
		if cellA == "" {
			break
		}

		cellB, _ := f.GetCellValue("Sheet1", fmt.Sprintf("B%d", i))
		if cellA == "Equal To" {
			value, _ := strconv.Atoi(cellB)
			results = append(results, value) // Equal Toの場合、B列の値を直接追加
			continue
		}

		cellC, _ := f.GetCellValue("Sheet1", fmt.Sprintf("C%d", i))
		if cellA == "In Range" {
			start, _ := strconv.Atoi(cellB)
			end, _ := strconv.Atoi(cellC)
			for j := start; j <= end; j++ {
				results = append(results, j) // In Rangeの場合、範囲内の全数値を追加
			}
		}
	}

	// 結果をJ列に出力
	for i, value := range results {
		f.SetCellInt("Sheet1", fmt.Sprintf("J%d", i+1), value)
	}

	// 変更を保存
	if err := f.SaveAs(filePath); err != nil {
		fmt.Println(err)
	}
}
