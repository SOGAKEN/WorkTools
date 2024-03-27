package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run script.go <ExcelFilePath>")
	}

	filePath := os.Args[1]
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer func() {
		// Save the file after modifications
		if err := f.Save(); err != nil {
			log.Fatalf("Failed to save file: %v", err)
		}
	}()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		log.Fatalf("Failed to get rows: %v", err)
	}

	jColumnIndex := 10 // J column index
	for i, row := range rows {
		if len(row) < 3 { // Skip rows that don't have enough data
			continue
		}
		keyword := strings.ToLower(row[0]) // Convert keyword to lowercase
		startNum, errStart := strconv.Atoi(row[1])
		endNum, errEnd := strconv.Atoi(row[2])
		if errStart != nil || errEnd != nil {
			fmt.Printf("Skipping row %d due to conversion error\n", i+1)
			continue
		}

		switch keyword {
		case "in range":
			for j := startNum; j <= endNum; j++ {
				cell, _ := excelize.CoordinatesToCellName(jColumnIndex, i+2) // i+2 because Excel is 1-indexed and starts from row 2 for data
				f.SetCellValue("Sheet1", cell, j)
			}
		case "equal to":
			cell, _ := excelize.CoordinatesToCellName(jColumnIndex, i+2)
			f.SetCellValue("Sheet1", cell, startNum)
		}
	}
}
