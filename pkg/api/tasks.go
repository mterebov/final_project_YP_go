package api

import (
	"encoding/json"
	"final_project_yp_go/pkg/db"
	"net/http"
)


type TasksResp struct {
    Tasks []*db.Task `json:"tasks"`
}


func writeJson(w http.ResponseWriter, tasks TasksResp) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(tasks)
}


func tasksHandle(w http.ResponseWriter, r *http.Request) {
    tasks, err := db.Tasks(50) // в параметре максимальное количество записей
    if err != nil {
		_ = RespCreator(w, 500, err, 0)
        return
    }
    writeJson(w, TasksResp{Tasks: tasks})
}