package main

import (
	"flag"
	"fmt"

	"golang/aws"
	"golang/office"
)

func init() {
	fmt.Println("test2.go の init() 関数")
}

func main() {
	var runAwsCpuCSV bool
	var csvFilePath string

	flag.BoolVar(&runAwsCpuCSV, "aws-cpu-csv", false, "Run AwsCpuCSV function")
	flag.StringVar(&csvFilePath, "file", "", "Path to CSV file for AwsCpuCSV")
	flag.Parse()

	if runAwsCpuCSV {
		aws.AwsCpuCsv(csvFilePath)
		return
	}
	office.ExcelAddColumn()

}
