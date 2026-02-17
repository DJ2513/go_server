package todo

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type TodoList struct {
	Id        string    `json:"id"`
	Title     string    `json:"title"`
	Tasks     []Task    `json:"tasks"`
	CreatedAt time.Time `json:"created_at"`
}

type Task struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	Done        bool      `json:"done"`
	Lists_id    string    `json:"-"`
}

func NewTodoList(title string) TodoList {
	return TodoList{
		Id:        uuid.New().String(),
		Title:     title,
		Tasks:     make([]Task, 0, 8),
		CreatedAt: time.Now(),
	}
}

func (t *TodoList) AddTask(task Task) error {

	task.CreatedAt = time.Now()
	task.Id = uuid.New().String()

	if task.Name == "" {
		return fmt.Errorf("[ERROR]: Tasks need a name")
	}

	t.Tasks = append(t.Tasks, task)
	return nil
}

func (t *TodoList) DeleteTask(id string) error {
	for i := 0; i < len(t.Tasks); i++ {
		if t.Tasks[i].Id == id {
			t.Tasks[i] = Task{}
		}
	}
	return fmt.Errorf("[ERROR]: Could not delete task %s.", id)
}
