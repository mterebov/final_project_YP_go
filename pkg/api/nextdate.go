package api

import (
	"fmt"
	"strings"
	"strconv"
	"net/http"
	"time"
)

const timePattern = "20060102"

func afterNow(date, now time.Time) bool {
	nowFormatted, _ := time.Parse(timePattern, now.Format(timePattern))
	return date.After(nowFormatted)
}


func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	// Парсим дату
	date, err := time.Parse(timePattern, dstart)
	if err != nil {
		return "", fmt.Errorf("time parse error: %s", err)
	}

	// Разделяем правило на слайс подстрок и проводим проверку их колличества
	repeatParts := strings.Split(repeat, " ")
	if (len(repeatParts) < 2 || len(repeatParts) > 4) && repeatParts[0] != "y" {
		return "", fmt.Errorf("unsupported format: wrong count of args")
	}

	// Обрабатывем полученое правило
	switch repeatParts[0] {
	// Для d <кол-во дней>
	case "d":
		// Конвертируем дни в int64
		interval, err := strconv.Atoi(repeatParts[1])

		// Проверяем ошибки и лимит интервала
		if err != nil || interval > 400 {
			if interval > 400 {
				return "", fmt.Errorf("interval limit: 400 < %d", interval)
			}
			return "", fmt.Errorf("interval parse error: %s", err)
		}
		
		// Вычисляем следуюющую дату
		for {
			date = date.AddDate(0, 0, interval)
			if afterNow(date, now) {
				break
			}
		}
		
		return date.Format(timePattern), nil
	
	// Для y (ежегодное повторение)
	case "y":
		if len(repeatParts) > 1 {
			return "", fmt.Errorf("unsupported format")
		}
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}
		return date.Format(timePattern), nil
	
	case "w":
		return "", fmt.Errorf("unsupported format: support may appear soon")
	case "m":
		return "", fmt.Errorf("unsupported format: support may appear soon")
	default:
		return "", fmt.Errorf("unsupported format")
	}
}


func NextDateHandle(res http.ResponseWriter, req *http.Request) {
	nowString := req.FormValue("now")
	date := req.FormValue("date")
	repeat := req.FormValue("repeat")
	nowTime, err := time.Parse(timePattern, nowString)
	if err != nil || len(nowString) == 0 {
		nowTime = time.Now()
	}
	response, err := NextDate(nowTime, date, repeat)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprint(res, response)
}