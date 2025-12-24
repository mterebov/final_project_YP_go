package api

import (
	"net/http"
)

func Init(mux *http.ServeMux) {
	mux.HandleFunc("/api/nextdate", NextDateHandle)
	mux.HandleFunc("/api/task", taskHandle)
	mux.HandleFunc("/api/tasks", tasksHandle)
}


func taskHandle(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	// Обработка создания задачи
	case http.MethodPost:
		addTaskHandle(res, req)
	// Обработка получения задачи
	case http.MethodGet:
		getTaskHandle(res, req)
	// Обработка редактирования задачи
	case http.MethodPut:
		updateTaskHandle(res, req)
	default:
		http.Error(res, "Bad Request: unreacheble", http.StatusBadRequest)
	}
}