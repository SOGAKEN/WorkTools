package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

func main() {
	aInums, err := readCsv("a.csv")
	if err != nil {
		panic(err)
	}

	bInums, err := readCsv("b.csv")
	if err != nil {
		panic(err)
	}

	aUnique, bUnique, common := compareLists(aInums, bInums)

	fmt.Println("Unique to a.csv:", aUnique)
	fmt.Println("Unique to b.csv:", bUnique)
	fmt.Println("Common:", common)
}

// readCsv reads the CSV file and returns a slice of inum values
func readCsv(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var inums []string
	for _, record := range records {
		inums = append(inums, record[0]) // Assuming 'inum' is the first column
	}

	return inums, nil
}

// compareLists compares two slices and returns unique and common elements
func compareLists(a, b []string) (aUnique, bUnique, common []string) {
	mA := make(map[string]bool)
	mB := make(map[string]bool)

	for _, item := range a {
		mA[item] = true
	}
	for _, item := range b {
		mB[item] = true
		if _, found := mA[item]; found {
			common = append(common, item)
		} else {
			bUnique = append(bUnique, item)
		}
	}
	for _, item := range a {
		if _, found := mB[item]; !found {
			aUnique = append(aUnique, item)
		}
	}

	return
}
