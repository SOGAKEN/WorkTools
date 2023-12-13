
package main

import (
    "encoding/csv"
    "os"
)

func main() {
    // ファイルを読み込む
    protectSet := readCSVToSet("protect.csv")
    filterCSV("list.csv", protectSet)
    filterCSV("one.csv", protectSet)
}

func readCSVToSet(filename string) map[string]bool {
    file, err := os.Open(filename)
    if err != nil {
        panic(err)
    }
    defer file.Close()

    lines, err := csv.NewReader(file).ReadAll()
    if err != nil {
        panic(err)
    }

    resultSet := make(map[string]bool)
    for _, line := range lines {
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

    lines, err := csv.NewReader(file).ReadAll()
    if err != nil {
        panic(err)
    }

    filtered := [][]string{}
    for _, line := range lines {
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
