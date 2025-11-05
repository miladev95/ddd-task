package main

import (
	"fmt"
	"log"
	"net/http"

	httpServer "github.com/example/task-management/interface/http"
	"github.com/example/task-management/shared/di"
)

func main() {
	// Initialize DI container
	container := di.NewContainer()

	// Create and setup HTTP router
	router := httpServer.NewRouter(container)
	router.SetupRoutes()

	// Start HTTP server
	port := ":8080"
	fmt.Printf("Starting Task Management API server on %s\n", port)

	if err := http.ListenAndServe(port, router.Handler()); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}