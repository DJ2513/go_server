// Package storage is in change of the database logic and comunication
package storage

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/DJ2513/go_server.git/internal/todo"
	"github.com/google/uuid"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) CreateTodoList(title string) (todo.TodoList, error) {
	td := todo.NewTodoList(title)
	_, err := r.DB.Exec(`
		INSERT INTO todo_lists 
		(id, title,created_at) 
		VALUES (? ,?, ?)
		`, td.ID, td.Title, td.CreatedAt)
	if err != nil {
		return todo.TodoList{}, fmt.Errorf("error insterting the list into de db:%v", err)
	}
	return td, nil
}

func (r *Repository) GetTodoList(id string) (todo.TodoList, error) {
	var td todo.TodoList
	row := r.DB.QueryRow(`
		SELECT id, title, created_at FROM todo_lists as td 
		WHERE td.id = ?`, id)

	err := row.Scan(&td.ID, &td.Title, &td.CreatedAt)
	if err != nil {
		return todo.TodoList{}, fmt.Errorf("error getting todo list from DB:%v", err)
	}

	res, err := r.DB.Query(`
		SELECT id, name, description, done, created_at
		FROM tasks as t
		WHERE t.list_id = ?`, id)

	if err != nil {
		return todo.TodoList{}, fmt.Errorf("error getting tasks from DB:%v", err)
	}
	defer res.Close()

	for res.Next() {
		var t todo.Task
		err := res.Scan(&t.ID, &t.Name, &t.Description, &t.Done, &t.CreatedAt)
		if err != nil {
			return todo.TodoList{}, fmt.Errorf("error scanning task:%v", err)
		}
		td.Tasks = append(td.Tasks, t)
	}

	return td, nil
}

func (r *Repository) GetAllTodoList() (*[]todo.TodoList, error) {
	var todos []todo.TodoList
	td := make(map[string]*todo.TodoList)
	rows, err := r.DB.Query(`
		SELECT id, title, created_at FROM todo_lists
	`)
	if err != nil {
		return &[]todo.TodoList{}, fmt.Errorf("error while getting todo lists: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var t todo.TodoList
		err = rows.Scan(&t.ID, &t.Title, &t.CreatedAt)
		if err != nil {
			return &[]todo.TodoList{}, fmt.Errorf("error scanning todo list: %v", err)
		}
		todos = append(todos, t)
	}

	for i := range todos {
		td[todos[i].ID] = &todos[i]
	}

	rows, err = r.DB.Query(`
		SELECT id, name, description, done, created_at, list_id 
		FROM tasks
	`)
	if err != nil {
		return &[]todo.TodoList{}, fmt.Errorf("error while getting tasks: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var t todo.Task
		err := rows.Scan(&t.ID, &t.Name, &t.Description, &t.Done, &t.CreatedAt, &t.ListsID)
		if err != nil {
			return &[]todo.TodoList{}, fmt.Errorf("error scanning tasks: %v", err)
		}
		if _, ok := td[t.ListsID]; !ok {
			return &[]todo.TodoList{}, fmt.Errorf("error, there is a task that has no list!: %v", t)
		}
		td[t.ListsID].Tasks = append(td[t.ListsID].Tasks, t)
	}

	return &todos, nil
}

func (r *Repository) DeleteTodoList(id string) error {
	_, err := r.DB.Exec(`
		DELETE FROM tasks
		WHERE list_id = ?`, id)
	if err != nil {
		return fmt.Errorf("error while deleting the todo list tasks: %v", err)
	}
	res, err := r.DB.Exec(`
		DELETE FROM todo_lists
		WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("error while deleting the todo list: %v", err)
	}
	x, _ := res.RowsAffected()
	if x == 0 {
		return fmt.Errorf("the list is none existant")
	}

	return nil
}

func (r *Repository) AddTask(listID string, task todo.Task) (todo.Task, error) {
	var exists string
	err := r.DB.QueryRow(`
		SELECT id FROM todo_lists
		WHERE id = ?`, listID).Scan(&exists)

	if err == sql.ErrNoRows {
		return todo.Task{}, fmt.Errorf("the tasklist %v does not exist", listID)
	}
	task.ID = uuid.New().String()
	task.CreatedAt = time.Now()

	_, err = r.DB.Exec(`
		INSERT INTO tasks 
		(id, name, description, done, created_at, list_id) 
		VALUES (? ,?, ?, ?, ?, ?)`,
		task.ID, task.Name, task.Description, task.Done, task.CreatedAt, listID,
	)
	if err != nil {
		return todo.Task{}, fmt.Errorf("error adding task to db: %v", err)
	}

	return task, nil
}

func (r *Repository) DeleteTask(taskID string) error {
	_, err := r.DB.Exec(`
		DELETE FROM tasks
		WHERE id = ?`, taskID)
	if err != nil {
		return fmt.Errorf("error while deleting the tasks: %v", err)
	}
	return nil
}
