package api

import (
	"fmt"
	"strings"
	"strconv"
	"net/http"
	"time"
)


func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	// Парсим дату
	date, err := time.Parse(TimePattern, dstart)
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
		
		return date.Format(TimePattern), nil
	
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
		return date.Format(TimePattern), nil
	
	case "w":
		if len(repeatParts) != 2 {
			return "", fmt.Errorf("unsupported format: w requires days list")
		}

		list, err := parseCSVInts(repeatParts[1])
		if err != nil {
			return "", fmt.Errorf("bad w list: %w", err)
		}

		var allowed [8]bool // 1..7
		for _, d := range list {
			if d < 1 || d > 7 {
				return "", fmt.Errorf("bad weekday: %d", d)
			}
			allowed[d] = true
		}

		// Стартуем с даты dstart и двигаемся до тех пор, пока дата не станет > now
		// А потом ищем ближайший подходящий weekday, день за днём.
		for {
			// перевести time.Weekday (Sun=0..Sat=6) в 1..7 (Mon=1..Sun=7)
			wd := int(date.Weekday())
			var wnum int
			if wd == 0 {
				wnum = 7
			} else {
				wnum = wd
			}

			if allowed[wnum] && afterNow(date, now) {
				return date.Format(TimePattern), nil
			}

			date = date.AddDate(0, 0, 1)
			// защитный лимит
			if date.Sub(now) > (time.Hour * 24 * 366 * 10) {
				return "", fmt.Errorf("cannot find next w date")
			}
		}
	case "m":
		if len(repeatParts) != 2 && len(repeatParts) != 3 {
			return "", fmt.Errorf("unsupported format: m requires days list and optional months list")
		}

		// дни месяца
		dayList, err := parseCSVInts(repeatParts[1])
		if err != nil {
			return "", fmt.Errorf("bad m days: %w", err)
		}

		var dayAllowed [32]bool // 1..31
		wantLast := false
		wantPrevLast := false
		for _, d := range dayList {
			switch d {
			case -1:
				wantLast = true
			case -2:
				wantPrevLast = true
			default:
				if d < 1 || d > 31 {
					return "", fmt.Errorf("bad day of month: %d", d)
				}
				dayAllowed[d] = true
			}
		}

		// месяцы
		var monthAllowed [13]bool // 1..12
		if len(repeatParts) == 3 {
			monthList, err := parseCSVInts(repeatParts[2])
			if err != nil {
				return "", fmt.Errorf("bad m months: %w", err)
			}
			for _, m := range monthList {
				if m < 1 || m > 12 {
					return "", fmt.Errorf("bad month: %d", m)
				}
				monthAllowed[m] = true
			}
		} else {
			for m := 1; m <= 12; m++ {
				monthAllowed[m] = true
			}
		}

		// идём по дням вперёд от dstart и ищем первый день > now, который подходит
		for {
			if afterNow(date, now) {
				mon := int(date.Month())
				if monthAllowed[mon] {
					d := date.Day()
					dim := daysInMonth(date)

					ok := dayAllowed[d]
					if !ok && wantLast && d == dim {
						ok = true
					}
					if !ok && wantPrevLast && d == dim-1 {
						ok = true
					}

					if ok {
						return date.Format(TimePattern), nil
					}
				}
			}

			date = date.AddDate(0, 0, 1)

			// защитный лимит
			if date.Sub(now) > (time.Hour * 24 * 366 * 20) {
				return "", fmt.Errorf("cannot find next m date")
			}
		}
	default:
		return "", fmt.Errorf("unsupported format")
	}
}


func NextDateHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	nowString := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	nowTime, err := time.Parse(TimePattern, nowString)
	if err != nil || len(nowString) == 0 {
		nowTime = time.Now()
	}

	response, err := NextDate(nowTime, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	fmt.Fprint(w, response)
}