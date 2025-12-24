package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"final_project_yp_go/pkg/db"
)


type response struct {
	ID 		int64 	`json:"id,omitempty"`
	Error 	string  `json:"error,omitempty"`
}


type TasksResp struct {
    Tasks []*db.Task `json:"tasks"`
}


const TimePattern = "20060102"


func afterNow(date, now time.Time) bool {
	y1, m1, d1 := date.Date()
	y2, m2, d2 := now.Date()
	return y1 > y2 || (y1 == y2 && m1 > m2) || (y1 == y2 && m1 == m2 && d1 > d2)
}


func daysInMonth(t time.Time) int {
	// последний день месяца: day=0 следующего месяца
	last := time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, t.Location())
	return last.Day()
}


func parseCSVInts(s string) ([]int, error) {
	if strings.TrimSpace(s) == "" {
		return nil, fmt.Errorf("empty list")
	}
	parts := strings.Split(s, ",")
	out := make([]int, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			return nil, fmt.Errorf("bad list")
		}
		n, err := strconv.Atoi(p)
		if err != nil {
			return nil, fmt.Errorf("bad number %q", p)
		}
		out = append(out, n)
	}
	return out, nil
}


func checkDate(task *db.Task) error {
	now := time.Now()

	// Подставить сегодняшнюю дату, если не указана
	if task.Date == "" {
		task.Date = now.Format(TimePattern)
	}

	// Проверка формата даты
	startDate, err := time.Parse(TimePattern, task.Date)
	if err != nil {
		return fmt.Errorf("Incorrect date format: %w", err)
	}

	// Проверить и пересчитать повтор
	if task.Repeat != "" {
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("Repeat rule error: %w", err)
		}
		// Если дата в прошлом — заменить на next
		if afterNow(now, startDate) {
			task.Date = next
		}
	} else {
		// Если повтора нет, но дата в прошлом — заменим на сегодня
		if afterNow(now, startDate) {
			task.Date = now.Format(TimePattern)
		}
	}

	return nil
}


func writeJson(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(v)
}