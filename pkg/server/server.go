package server

import (
	"net/http"
	"os"
	"strings"
	"fmt"

	"final_project_yp_go/pkg/db"
	"final_project_yp_go/pkg/api"
)

func StartServer() error {
	dbFile := strings.TrimSpace(os.Getenv("TODO_DBFILE"))
	fmt.Println(os.Getenv("TODO_DBFILE"))
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	if err := db.Init(dbFile); err != nil {
		return fmt.Errorf("db init error: %w", err)
	}
	defer db.Close()

	port := strings.TrimSpace(os.Getenv("TODO_PORT"))
	if port == "" {
		port = "7540"
	}

	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}
	fmt.Printf("Server port: %s\nDatabase file: %s\n", port, dbFile)
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./web")))
	api.Init(mux)
	if err := http.ListenAndServe(port, mux); err != nil {
		return fmt.Errorf("server start error: %s", err)
	}
	return nil
}
