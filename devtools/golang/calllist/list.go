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

	// J列に出力する際の開始行を追跡
	jColumnRow := 1

	for i := 1; ; i++ { // A列が空になるまで繰り返し
		cellA, _ := f.GetCellValue("Sheet1", fmt.Sprintf("A%d", i))
		if cellA == "" {
			break
		}

		cellB, _ := f.GetCellValue("Sheet1", fmt.Sprintf("B%d", i))
		cellC, _ := f.GetCellValue("Sheet1", fmt.Sprintf("C%d", i))

		switch cellA {
		case "In Range":
			start, _ := strconv.Atoi(cellB)
			end, _ := strconv.Atoi(cellC)
			for j := start; j <= end; j++ {
				f.SetCellInt("Sheet1", fmt.Sprintf("J%d", jColumnRow), j)
				jColumnRow++
			}
		case "Equal To":
			value, _ := strconv.Atoi(cellB)
			f.SetCellInt("Sheet1", fmt.Sprintf("J%d", jColumnRow), value)
			jColumnRow++
		}
	}

	// 変更を保存
	if err := f.SaveAs(filePath); err != nil {
		fmt.Println(err)
	}
}
