package types

import "time"

type Todo struct {
	ID          string    `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Category    string    `json:"category" db:"category"`
	IsCompleted bool      `json:"is_completed" db:"is_completed"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	UserID      string    `json:"user_id" db:"user_id"`
}

type TodoPayload struct {
	Title       string `json:"title" required:"true"`
	Description string `json:"description" required:"true"`
	Category    string `json:"category" required:"true"`
	IsConpleted bool   `json:"is_completed" required:"true"`
	UserId      string `json:"user_id" required:"true"`
}

type TodoStore interface {
	CreateTodo(payload TodoPayload) error
	DeleteTodo(id string, userId string) error
	UpdateTodo(id string, userId string, payload TodoPayload) error
	SetIsCompledtedTodo(id string, userId string) error
	GetTodoById(id string) (*Todo, error)
	GetTodoByDate(date string, id string) (*[]Todo, error)
}
