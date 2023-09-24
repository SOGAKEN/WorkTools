package aws

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func AwsCpuCsv(csvFile string) {

	f, err := os.Open(csvFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = -1

	// 1-3行目をスキップ
	for i := 0; i < 3; i++ {
		_, err := r.Read()
		if err != nil {
			log.Fatal(err)
		}
	}

	// 4行目をヘッダーとして読み込み
	headers, err := r.Read()
	if err != nil {
		log.Fatal(err)
	}

	// ヘッダーの修正（指定された文字列を含む場合に適用）
	headerMapping := map[string]string{
		"i-039f02b35a9336a84": "Chat",
		"i-039f02b35a9336a82": "Chat2",
	}

	for idx, header := range headers {
		for key, newValue := range headerMapping {
			if strings.Contains(header, key) {
				headers[idx] = newValue
				break
			}
		}
	}

	data := [][]string{{"Label", "Average", "Max", "Date of Max"}}
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// 列ごとの計算
	for col := 1; col < len(headers); col++ {
		var sum float64
		var count int
		var max float64 = -1e9
		var dateOfMax string

		for _, record := range records {
			if len(record) <= col {
				continue
			}

			if v, err := strconv.ParseFloat(record[col], 64); err == nil {
				sum += v
				count++
				if v > max {
					max = v
					dateOfMax = record[0]
				}
			}
		}

		if count == 0 {
			continue
		}
		avg := sum / float64(count)

		data = append(data, []string{headers[col], fmt.Sprintf("%f", avg), fmt.Sprintf("%f", max), dateOfMax})
	}

	// 日付取得
	now := time.Now()
	date := now.Format("20060102")

	// 出力ファイル作成
	outfile := "new_" + date + ".csv"
	outf, err := os.Create(outfile)
	if err != nil {
		log.Fatal(err)
	}
	defer outf.Close()

	// CSV出力
	w := csv.NewWriter(outf)
	if err := w.WriteAll(data); err != nil {
		log.Fatal(err)
	}

	w.Flush()
	fmt.Println("Done!")
}
