package db

import (
	"time"

)


func searchToDate(search string) (string, error) {
	date, err := time.Parse("02.01.2006", search)
	if err != nil {
		return "", err
	}
	return date.Format("20060102"), nil
}