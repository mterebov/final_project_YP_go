package api

import (
	"net/http"
)

func Init(mux *http.ServeMux) {
	mux.HandleFunc("/api/nextdate", NextDateHandle)
	mux.HandleFunc("/api/task", taskHandle)
}


func taskHandle(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	// Обработка создания задачи
	case http.MethodPost:
		addTaskHandle(res, req)
	// Обработка получения задачи
	// case http.MethodGet:
	
	// Обработка удаления задачи
	// case http.MethodDelete:
	
	default:
		http.Error(res, "Bad Request: unreacheble", http.StatusBadRequest)
	}
}