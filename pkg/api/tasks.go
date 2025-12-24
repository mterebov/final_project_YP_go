package api

import (
	"encoding/json"
	"final_project_yp_go/pkg/db"
	"net/http"
)


type TasksResp struct {
    Tasks []*db.Task `json:"tasks"`
}


func writeJson(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(v)
}


func tasksHandle(w http.ResponseWriter, r *http.Request) {
    tasks, err := db.Tasks(50) // в параметре максимальное количество записей
    if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
        writeJson(w, response{Error: err.Error()})
        return
    }
	w.WriteHeader(http.StatusOK)
    writeJson(w, TasksResp{Tasks: tasks})
}