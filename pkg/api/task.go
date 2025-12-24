package api

import (
	"database/sql"
	"encoding/json"
	"final_project_yp_go/pkg/db"
	"net/http"
	"strings"
	"time"
)


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


func doneTaskHandle(w http.ResponseWriter, r *http.Request) {
    id := strings.TrimSpace(r.FormValue("id"))
    if id == "" {
        w.WriteHeader(http.StatusBadRequest)
        writeJson(w, response{Error: "id is required"})
        return
    }
    task, err := db.GetTask(id)
    if err != nil {
        if err == sql.ErrNoRows {
            w.WriteHeader(http.StatusNotFound)
        } else {
            w.WriteHeader(http.StatusInternalServerError)
        }
        writeJson(w, response{Error: err.Error()})
        return
    }
    if task.Repeat == "" {
        err := db.DeleteTask(task.ID)
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            writeJson(w, response{Error: err.Error()})
            return
        }
        w.WriteHeader(http.StatusOK)
        writeJson(w, response{})
        return
    }
    next, err := NextDate(time.Now(), task.Date, task.Repeat)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        writeJson(w, response{Error: err.Error()})
        return
    }
    task.Date = next
    if err = db.UpdateTask(task); err != nil {
        if err == sql.ErrNoRows {
            w.WriteHeader(http.StatusNotFound)
        } else {
            w.WriteHeader(http.StatusInternalServerError)
        }
        writeJson(w, response{Error: err.Error()})
        return
    }
    w.WriteHeader(http.StatusOK)
    writeJson(w, response{})
}


func deleteTaskHandle(w http.ResponseWriter, r *http.Request) {
    id := strings.TrimSpace(r.FormValue("id"))
    if id == "" {
        w.WriteHeader(http.StatusBadRequest)
        writeJson(w, response{Error: "id is required"})
        return
    }
    err := db.DeleteTask(id)
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
    writeJson(w, response{})
}
