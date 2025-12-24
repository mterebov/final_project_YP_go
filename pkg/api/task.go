package api

import (
	"database/sql"
	"encoding/json"
	"final_project_yp_go/pkg/db"
	"fmt"
	"net/http"
	"strings"
	"time"
)


type response struct {
	ID 		int64 `json:"id,omitempty"`
	Error 	string  `json:"error,omitempty"`
}


func checkDate(task *db.Task) error {
	now := time.Now()

	// Подставить сегодняшнюю дату, если не указана
	if task.Date == "" {
		task.Date = now.Format("20060102")
	}

	// Проверка формата даты
	startDate, err := time.Parse("20060102", task.Date)
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
			task.Date = now.Format("20060102")
		}
	}

	return nil
}


func addTaskHandle(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    var task db.Task
    dec := json.NewDecoder(r.Body)

    if err := dec.Decode(&task); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        writeJson(w, response{Error: err.Error()})
        return
    }

    if strings.TrimSpace(task.Title) == "" {
        w.WriteHeader(http.StatusBadRequest)
        writeJson(w, response{Error: "title can't be empty"})
        return
    }

    if err := checkDate(&task); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        writeJson(w, response{Error: err.Error()})
        return
    }

    id, err := db.AddTask(&task)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        writeJson(w, response{Error: err.Error()})
        return
    }

    w.WriteHeader(http.StatusOK)
    writeJson(w, response{ID: id})
}


func getTaskHandle(w http.ResponseWriter, r *http.Request) {
    id := strings.TrimSpace(r.FormValue("id"))
    if id == "" {
        w.WriteHeader(http.StatusBadRequest)
        writeJson(w, response{Error: "id is required"})
        return
    }
    taskById, err := db.GetTask(id)
    if err != nil {
        if err == sql.ErrNoRows {
            w.WriteHeader(http.StatusNotFound)
        } else {
            w.WriteHeader(http.StatusInternalServerError)
        }
        writeJson(w, response{Error: err.Error()})
        return
    }
    w.WriteHeader(http.StatusOK)
    writeJson(w, taskById)
}


func updateTaskHandle(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    var task db.Task
    dec := json.NewDecoder(r.Body)

    if err := dec.Decode(&task); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        writeJson(w, response{Error: err.Error()})
        return
    }

    if strings.TrimSpace(task.Title) == "" {
        w.WriteHeader(http.StatusBadRequest)
        writeJson(w, response{Error: "title can't be empty"})
        return
    }

    if err := checkDate(&task); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        writeJson(w, response{Error: err.Error()})
        return
    }

    err := db.UpdateTask(&task)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        writeJson(w, response{Error: err.Error()})
        return
    }

    w.WriteHeader(http.StatusOK)
    writeJson(w, response{})
}