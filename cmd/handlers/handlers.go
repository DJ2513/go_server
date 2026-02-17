// package handlers

// Here we have the connection of the routes and the logic
// something like middleware witout being middleware

package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/DJ2513/go_server.git/internal/storage"
	"github.com/DJ2513/go_server.git/internal/todo"
)

type Handler struct {
	Repo *storage.Repository
}

func (h *Handler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"message":"Everything is good and running", "status":"%v", "timestamp":"%s"}`, http.StatusOK, time.Now())))
}

func (h *Handler) CreateTodoList(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string `json:"title"`
	}
	json.NewDecoder(r.Body).Decode(&input)
	todo, err := h.Repo.CreateTodoList(input.Title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func (h *Handler) GetTodoList(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	todo, err := h.Repo.GetTodoList(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func (h *Handler) GetAllTodoList(w http.ResponseWriter, _ *http.Request) {
	todos, err := h.Repo.GetAllTodoList()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func (h *Handler) DeleteTodoList(w http.ResponseWriter, r *http.Request) {
	todo_id := r.PathValue("id")
	err := h.Repo.DeleteTodoList(todo_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent) // 204
}

func (h *Handler) AddTask(w http.ResponseWriter, r *http.Request) {
	var task todo.Task
	list_id := r.PathValue("id")
	json.NewDecoder(r.Body).Decode(&task)
	created_task, err := h.Repo.AddTask(list_id, task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created_task)
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	task_id := r.PathValue("taskId")
	err := h.Repo.DeleteTask(task_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent) // 204
}
