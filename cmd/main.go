package main

import (
	"fmt"
	"net/http"
	"os"

	h "github.com/DJ2513/go_server.git/cmd/handlers"
	"github.com/joho/godotenv"
)

func main() {
	// Env variables
	if err := godotenv.Load(); err != nil {
		panic(fmt.Sprintf("Error loading env variables: %v", err))
	}
	port := os.Getenv("PORT")

	// Handlers
	http.HandleFunc("GET /", h.Homehandler)
	http.HandleFunc("GET /health", h.HealthCheckHandler)

	// Server
	fmt.Printf("Server running! Listening on https://localhost:%v\n\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(fmt.Sprintf("Error starting server: %v", err))
	}
}
