package api

import (
	"net/http"
)

func Init(mux *http.ServeMux) {
	mux.HandleFunc("/api/nextdate", NextDateHandle)
	mux.HandleFunc("/api/task", addTaskHandle)
	mux.HandleFunc("/api/tasks", tasksHandle)
}
