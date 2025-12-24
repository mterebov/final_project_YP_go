package api

import (
	"final_project_yp_go/pkg/db"
	"net/http"
)


func tasksHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	search := r.FormValue("search")

    tasks, err := db.Tasks(50, search) // в параметре максимальное количество записей
    if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
        writeJson(w, response{Error: err.Error()})
        return
    }
	w.WriteHeader(http.StatusOK)
    writeJson(w, TasksResp{Tasks: tasks})
}