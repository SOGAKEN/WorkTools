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
		cellA, errA := f.GetCellValue("Sheet1", fmt.Sprintf("A%d", i))
		if errA != nil || cellA == "" {
			break // エラーまたはA列が空なら処理を終了
		}

		cellB, errB := f.GetCellValue("Sheet1", fmt.Sprintf("B%d", i))
		if errB != nil {
			continue // B列の読み取りに失敗した場合はスキップ
		}

		switch cellA {
		case "In Range":
			start, errStart := strconv.Atoi(cellB)
			cellC, _ := f.GetCellValue("Sheet1", fmt.Sprintf("C%d", i))
			end, errEnd := strconv.Atoi(cellC)
			if errStart == nil && errEnd == nil {
				for j := start; j <= end; j++ {
					results = append(results, j) // In Rangeの場合、範囲内の全数値を追加
				}
			}
		case "Equal To":
			if value, err := strconv.Atoi(cellB); err == nil {
				results = append(results, value) // Equal Toの場合、B列の値を直接追加
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
