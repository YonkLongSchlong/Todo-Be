package types

import "time"

type User struct {
	ID        string    `json:"id" db:"id"`
	FirstName string    `json:"first_name" db:"first_name"`
	LastName  string    `json:"last_name" db:"last_name"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"password" db:"password"`
	Avatar    string    `json:"avatar" db:"avatar"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type UserPayload struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type PasswordPayload struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type UserStore interface {
	CreateUser(payload RegisterPayload) error
	GetUserByEmail(email string) (*User, error)
	GetUserById(id string) (*User, error)
	UpdateUser(id string, payload UserPayload) (*User, error)
	UpdatePassword(id string, hashesPassword string) error
	UpdateAvatar(id string, imageUrl string) (*User, error)
}
