package api

import (
	"encoding/json"
	"errors"
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
	var next string
	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format("20060102")
	}
	t, err := time.Parse("20060102", task.Date)
	if err != nil {return err}
	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {return err}
	}
	if afterNow(now, t) {
        if len(task.Repeat) == 0 {
            task.Date = now.Format("20060102")
        } else {
            task.Date = next
        }
    }
	return nil
}


func respCreator(w http.ResponseWriter, status int, err error, id int64) error {
    w.Header().Set("Content-Type", "application/json; charset=utf-8")

    var resBody response
    if err != nil {
        resBody.Error = err.Error()
        w.WriteHeader(status)
        return json.NewEncoder(w).Encode(resBody)
    }

    resBody.ID = id
    w.WriteHeader(status)
    return json.NewEncoder(w).Encode(resBody)
}


func addTaskHandle(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()

    var task db.Task
    dec := json.NewDecoder(r.Body)

    if err := dec.Decode(&task); err != nil {
        _ = respCreator(w, http.StatusBadRequest, err, 0)
        return
    }

    if strings.TrimSpace(task.Title) == "" {
        _ = respCreator(w, http.StatusBadRequest, errors.New("title is required"), 0)
        return
    }

    if err := checkDate(&task); err != nil {
        _ = respCreator(w, http.StatusBadRequest, err, 0)
        return
    }

    id, err := db.AddTask(&task)
    if err != nil {
		fmt.Println(err)
        _ = respCreator(w, http.StatusInternalServerError, err, 0)
        return
    }

    _ = respCreator(w, http.StatusOK, nil, id)
}

