package todo

import (
	"fmt"
	"time"

	"github.com/YonkLongSchlong/Todo-BE/packages/types"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateTodo(payload types.TodoPayload) error {
	todo := types.Todo{
		ID:          uuid.NewString(),
		Title:       payload.Title,
		Description: payload.Description,
		Category:    payload.Category,
		IsCompleted: payload.IsConpleted,
		CreatedAt:   time.Now().Local().Add(time.Hour * 7), // THIS IS FOR VIETNAME LOCALTIME TO SWITCH TO UTC SIMPLY USE TIME.NOW()
		UpdatedAt:   time.Now().Local().Add(time.Hour * 7), // THIS IS FOR VIETNAME LOCALTIME TO SWITCH TO UTC SIMPLY USE TIME.NOW()
		UserID:      payload.UserId,
	}

	fmt.Println(time.Now())
	_, err := s.db.NamedExec("INSERT INTO todos(id, title, description, category, is_completed, created_at, updated_at, user_id) VALUES (:id, :title, :description, :category, :is_completed, :created_at, :updated_at, :user_id)", todo)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) DeleteTodo(id string, userId string) error {
	/** CHECK IF THE TODO USER ID MATCHES WITH THE ID FROM REQUEST
	 * If true continue delete todo
	 * Else return error
	 */
	todo := new(types.Todo)
	err := s.db.Get(todo, "SELECT * FROM todos WHERE id = ? AND user_id = ?", id, userId)
	if err != nil {
		return fmt.Errorf("you don't have permission to delete this todo")
	}

	/** DELETE TODO */
	_, err = s.db.Exec("DELETE FROM todos WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) UpdateTodo(id string, userId string, payload types.TodoPayload) error {
	/** CHECK IF THE TODO USER ID MATCHES WITH THE ID FROM REQUEST
	 * If true continue delete todo
	 * Else return error
	 */
	todo := new(types.Todo)
	err := s.db.Get(todo, "SELECT * FROM todos WHERE id = ? AND user_id = ?", id, userId)
	if err != nil {
		return fmt.Errorf("you don't have permission to update this todo")
	}

	/** UPDATE TODO */
	_, err = s.db.Exec("UPDATE todos SET title = ?, description = ?, category = ?, updated_at = ? WHERE id = ?", payload.Title, payload.Description, payload.Category, time.Now().Local().Add(time.Hour*7), todo.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) SetIsCompledtedTodo(id string, userId string) error {
	/** CHECK IF THE TODO USER ID MATCHES WITH THE ID FROM REQUEST
	 * If true continue delete todo
	 * Else return error
	 */
	todo := new(types.Todo)
	err := s.db.Get(todo, "SELECT * FROM todos WHERE id = ? AND user_id = ?", id, userId)
	if err != nil {
		return fmt.Errorf("you don't have permission to update this todo")
	}

	/** UPDATE IS_COMPLETED  */
	_, err = s.db.Exec("UPDATE todos SET is_completed = ? WHERE id = ?", !todo.IsCompleted, todo.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetTodoById(id string) (*types.Todo, error) {
	todo := new(types.Todo)
	err := s.db.Get(todo, "SELECT * FROM todos WHERE id = ? ", id)
	if err != nil {
		return nil, fmt.Errorf("you don't have permission to update this todo")
	}
	return todo, nil
}

func (s *Store) GetTodoByDate(date string) (*[]types.Todo, error) {
	dateStart := fmt.Sprint(date + " 00:00:00")
	dateEnd := fmt.Sprint(date + " 23:59:00")
	todo := new([]types.Todo)
	err := s.db.Select(todo, "SELECT * FROM todos WHERE created_at BETWEEN ? AND ? ", dateStart, dateEnd)
	if err != nil {
		return nil, err
	}

	return todo, nil
}
