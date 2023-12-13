
package main

import (
    "encoding/csv"
    "os"
)

func main() {
    // protect.csvから値を読み込む
    protectSet := readCSVToSet("protect.csv")

    // list.csvとone.csvをフィルタリング
    filterCSV("list.csv", protectSet)
    filterCSV("one.csv", protectSet)
}

func readCSVToSet(filename string) map[string]bool {
    file, err := os.Open(filename)
    if err != nil {
        panic(err)
    }
    defer file.Close()

    resultSet := make(map[string]bool)
    csvReader := csv.NewReader(file)
    for {
        line, err := csvReader.Read()
        if err != nil {
            if err == csv.ErrFieldCount || err == io.EOF {
                break
            }
            panic(err)
        }
        resultSet[line[0]] = true
    }

    return resultSet
}

func filterCSV(filename string, protectSet map[string]bool) {
    file, err := os.Open(filename)
    if err != nil {
        panic(err)
    }
    defer file.Close()

    filtered := [][]string{}
    csvReader := csv.NewReader(file)
    for {
        line, err := csvReader.Read()
        if err != nil {
            if err == csv.ErrFieldCount || err == io.EOF {
                break
            }
            panic(err)
        }
        if _, found := protectSet[line[0]]; !found {
            filtered = append(filtered, line)
        }
    }

    writeFile(filename, filtered)
}

func writeFile(filename string, data [][]string) {
    file, err := os.Create(filename)
    if err != nil {
        panic(err)
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    for _, line := range data {
        if err := writer.Write(line); err != nil {
            panic(err)
        }
    }
}
