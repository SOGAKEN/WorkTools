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
	SecondLineValue  int //対象期間
}

type SecondQueryParams struct {
	DB              *sql.DB
	OneWriter       *csv.Writer
	ProtectedValues map[string]struct{}
}
