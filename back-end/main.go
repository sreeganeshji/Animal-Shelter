package main

import (
	"back-end/database"
	"back-end/rest"
	"back-end/service"
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Starting application...")

	// Create database object
	postgres, err := database.Connect("localhost", 5432, "postgres", "postgres", "postgres", 3)
	if err != nil {
		panic(err)
	}

	// Close Postgres connection if main terminates
	defer postgres.Close()

	// Create service object
	service := service.Service{Postgres: postgres}

	// Create REST routes
	r := rest.Init(service)

	// Listen on port
	log.Fatal(http.ListenAndServe("localhost:8080", r))
}
