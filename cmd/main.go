package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/DJ2513/go_server.git/cmd/handlers"
	"github.com/DJ2513/go_server.git/internal/storage"
	"github.com/joho/godotenv"
)

func main() {
	// Env variables
	if err := godotenv.Load(); err != nil {
		panic(fmt.Sprintf("Error loading env variables: %v", err))
	}
	port := os.Getenv("PORT")

	db, err := storage.InitDB("./todo.db")
	if err != nil {
		panic(fmt.Sprintf("Error initializing db: %v", err))
	}
	defer db.Close()

	// Handlers
	h := handlers.Handler{
		Repo: storage.NewRepository(db),
	}

	http.HandleFunc("GET /health", h.HealthCheckHandler)
	http.HandleFunc("POST /todos", h.CreateTodoList)
	http.HandleFunc("GET /todos", h.GetAllTodoList)
	http.HandleFunc("GET /todos/{id}", h.GetTodoList)
	http.HandleFunc("DELETE /todos/{id}", h.DeleteTodoList)
	http.HandleFunc("POST /todos/{id}/tasks", h.AddTask)
	http.HandleFunc("DELETE /todos/{id}/tasks/{taskId}", h.DeleteTask)

	// Server
	fmt.Printf("Server running! Listening on http://localhost:%v\n\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(fmt.Sprintf("Error starting server: %v", err))
	}
}
