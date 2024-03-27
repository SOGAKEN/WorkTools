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
		fmt.Println("Error opening file:", err)
		return
	}

	var results []int // 結果を格納するためのスライス

	for i := 1; ; i++ { // A列が空になるまで繰り返し
		cellA, err := f.GetCellValue("Sheet1", fmt.Sprintf("A%d", i))
		if err != nil || cellA == "" {
			break // エラーまたはA列が空なら処理を終了
		}

		cellB, _ := f.GetCellValue("Sheet1", fmt.Sprintf("B%d", i))

		switch cellA {
		case "In Range":
			start, errStart := strconv.Atoi(cellB)
			if errStart != nil {
				fmt.Printf("Error converting B%d: %v\n", i, errStart)
				continue
			}
			cellC, _ := f.GetCellValue("Sheet1", fmt.Sprintf("C%d", i))
			end, errEnd := strconv.Atoi(cellC)
			if errEnd != nil {
				fmt.Printf("Error converting C%d: %v\n", i, errEnd)
				continue
			}
			for j := start; j <= end; j++ {
				results = append(results, j)
			}
		case "Equal To":
			value, err := strconv.Atoi(cellB)
			if err != nil {
				fmt.Printf("Error converting B%d: %v\n", i, err)
				continue // 変換エラーがあればこの行をスキップ
			}
			results = append(results, value)
		}
	}

	// 結果をJ列に出力
	for i, value := range results {
		f.SetCellInt("Sheet1", fmt.Sprintf("J%d", i+1), value)
	}

	// 変更を保存
	if err := f.SaveAs(filePath); err != nil {
		fmt.Println("Error saving file:", err)
	}
}
