package model

import (
	"database/sql"
	"encoding/csv"
)

type DataRecord struct {
	LogID                string
	MinDate              string
	MaxDate              string
	HowManyDaysFromToday int
}

type FirstQueryParams struct {
	DB               *sql.DB
	AllWriter        *csv.Writer
	ListWriter       *csv.Writer
	NotProtectWriter *csv.Writer
	ProtectedValues  map[string]struct{}
	SecondLineValue  int
}
