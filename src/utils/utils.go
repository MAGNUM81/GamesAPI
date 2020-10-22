package utils

import "time"

func GetDate(strDate string) time.Time {
	const dateFormat = "2006-01-02"
	date, _ := time.Parse(dateFormat, strDate)
	return date
}
