package util

import (
	"time"
)

const (
	FORMAT_DATETIME = "2006-01-02 15:04:05"
	FORMAT_DATE     = "2006-01-02"
)

func Now() string {

	return time.Now().Format(FORMAT_DATETIME)
}

func NowDate() string {

	return time.Now().Format(FORMAT_DATE)
}

func Date(date string) string {

	if formatedDate, err := time.Parse(FORMAT_DATETIME, date); err != nil {
		return "0000-00-00"
	} else {
		return formatedDate.Format(FORMAT_DATE)
	}
}
