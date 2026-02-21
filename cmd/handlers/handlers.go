// Package handlers is the bridge between the user requests and the db requests.
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
	todoID := r.PathValue("id")
	err := h.Repo.DeleteTodoList(todoID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent) // 204
}

func (h *Handler) AddTask(w http.ResponseWriter, r *http.Request) {
	var task todo.Task
	listID := r.PathValue("id")
	json.NewDecoder(r.Body).Decode(&task)
	createTask, err := h.Repo.AddTask(listID, task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createTask)
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	taskID := r.PathValue("task_id")
	err := h.Repo.DeleteTask(taskID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent) // 204
}
