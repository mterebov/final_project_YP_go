package main

import (
	"log"

	"final_project_yp_go/pkg/server"
)


func main() {
	err := server.StartServer()
	if err != nil {
		log.Fatal(err)
	}
}