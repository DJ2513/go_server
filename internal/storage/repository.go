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
		`, td.Id, td.Title, td.CreatedAt)
	if err != nil {
		return todo.TodoList{}, fmt.Errorf("Error insterting the list into de db:%v", err)
	}
	return td, nil
}

func (r *Repository) GetTodoList(id string) (todo.TodoList, error) {
	var td todo.TodoList
	row := r.DB.QueryRow(`
		SELECT id, title, created_at FROM todo_lists as td 
		WHERE td.id = ?`, id)

	err := row.Scan(&td.Id, &td.Title, &td.CreatedAt)
	if err != nil {
		return todo.TodoList{}, fmt.Errorf("Error getting todo list from DB:%v", err)
	}

	res, err := r.DB.Query(`
		SELECT id, name, description, done, created_at
		FROM tasks as t
		WHERE t.list_id = ?`, id)

	if err != nil {
		return todo.TodoList{}, fmt.Errorf("Error getting tasks from DB:%v", err)
	}
	defer res.Close()

	for res.Next() {
		var t todo.Task
		err := res.Scan(&t.Id, &t.Name, &t.Description, &t.Done, &t.CreatedAt)
		if err != nil {
			return todo.TodoList{}, fmt.Errorf("Error scanning task:%v", err)
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
		return &[]todo.TodoList{}, fmt.Errorf("Error while getting todo lists: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var t todo.TodoList
		err := rows.Scan(&t.Id, &t.Title, &t.CreatedAt)
		if err != nil {
			return &[]todo.TodoList{}, fmt.Errorf("Error scanning todo list: %v", err)
		}
		todos = append(todos, t)
	}

	for i := range todos {
		td[todos[i].Id] = &todos[i]
	}

	rows, err = r.DB.Query(`
		SELECT id, name, description, done, created_at, list_id 
		FROM tasks
	`)
	if err != nil {
		return &[]todo.TodoList{}, fmt.Errorf("Error while getting tasks: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var t todo.Task
		err := rows.Scan(&t.Id, &t.Name, &t.Description, &t.Done, &t.CreatedAt, &t.Lists_id)
		if err != nil {
			return &[]todo.TodoList{}, fmt.Errorf("Error scanning tasks: %v", err)
		}
		if _, ok := td[t.Lists_id]; !ok {
			return &[]todo.TodoList{}, fmt.Errorf("Error, there is a task that has no list!: %v", t)
		}
		td[t.Lists_id].Tasks = append(td[t.Lists_id].Tasks, t)
	}

	return &todos, nil
}

func (r *Repository) DeleteTodoList(id string) error {
	_, err := r.DB.Exec(`
		DELETE FROM tasks
		WHERE list_id = ?`, id)
	if err != nil {
		return fmt.Errorf("Error while deleting the todo list tasks: %v", err)
	}
	res, err := r.DB.Exec(`
		DELETE FROM todo_lists
		WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("Error while deleting the todo list: %v", err)
	}
	x, _ := res.RowsAffected()
	if x == 0 {
		return fmt.Errorf("The list is none existant!")
	}

	return nil
}

func (r *Repository) AddTask(listID string, task todo.Task) (todo.Task, error) {
	var exists string
	err := r.DB.QueryRow(`
		SELECT id FROM todo_lists
		WHERE id = ?`, listID).Scan(&exists)

	if err == sql.ErrNoRows {
		return todo.Task{}, fmt.Errorf("The tasklist %v does not exist!", listID)
	}
	task.Id = uuid.New().String()
	task.CreatedAt = time.Now()

	_, err = r.DB.Exec(`
		INSERT INTO tasks 
		(id, name, description, done, created_at, list_id) 
		VALUES (? ,?, ?, ?, ?, ?)`,
		task.Id, task.Name, task.Description, task.Done, task.CreatedAt, listID,
	)
	if err != nil {
		return todo.Task{}, fmt.Errorf("Error adding task to db: %v", err)
	}

	return task, nil
}

func (r *Repository) DeleteTask(taskID string) error {
	_, err := r.DB.Exec(`
		DELETE FROM tasks
		WHERE id = ?`, taskID)
	if err != nil {
		return fmt.Errorf("Error while deleting the tasks: %v", err)
	}
	return nil
}
