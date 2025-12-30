package db

import (
	"time"

)

const TimePattern = "20060102"

func searchToDate(search string) (string, error) {
	date, err := time.Parse("02.01.2006", search)
	if err != nil {
		return "", err
	}
	// нужно использовать константу вместо 20060102
	return date.Format(TimePattern), nil
}