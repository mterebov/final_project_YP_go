package api

import (
	"net/http"
)

func Init(mux *http.ServeMux) {
	mux.HandleFunc("/api/nextdate", NextDateHandle)
	mux.HandleFunc("/api/task", taskHandle)
	mux.HandleFunc("/api/tasks", tasksHandle)
	mux.HandleFunc("/api/task/done", doneTaskHandle)
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
	// Обработка удаления задачи
	case http.MethodDelete:
		deleteTaskHandle(res, req)
	default:
		// тут код 405 должен быть
		http.Error(res, "Method Not Allowed", http.StatusMethodNotAllowed) // 405
	}
}