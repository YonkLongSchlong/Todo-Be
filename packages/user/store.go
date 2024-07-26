package user

import (
	"fmt"
	"time"

	"github.com/YonkLongSchlong/Todo-BE/packages/types"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateUser(payload types.RegisterPayload) error {
	/** CHECK IF THE USER WITH PAYLOAD EMAIL EXIST
	 * If exist return new error
	 * Else continue
	 */
	_, err := s.GetUserByEmail(payload.Email)
	if err == nil {
		return fmt.Errorf("this user with email %s already registered", payload.Email)
	}

	/** HASH THE PASSWORD */
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	/** CREATE USER AND INSERT TO DB
	 * Create a new user struct that hold payload values and hash password
	 * Insert to db if error return it
	 * Or el return nil
	 */
	user := types.User{
		ID:        uuid.New().String(),
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  string(hashPassword),
		Avatar:    "https://todo-avatar.s3.ap-southeast-1.amazonaws.com/placeholder.jpg",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err = s.db.NamedExec("INSERT INTO users(id, first_name, last_name, email, password, avatar, created_at, updated_at) VALUES (:id, :first_name, :last_name, :email,:password, :avatar, :created_at, :updated_at)", user)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	/** MAKE A NEW POINTER TO USER STRUCT */
	user := new(types.User)

	/** QUERY USER WITH EMAIL
	 * If user with email not exist return nil,err
	 * Else return user, nil
	 */
	err := s.db.Get(user, "SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %v", email)
	}
	return user, nil
}

func (s *Store) GetUserById(id string) (*types.User, error) {
	user := new(types.User)

	err := s.db.Get(user, "SELECT * from users WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Store) UpdateUser(id string, payload types.UserPayload) (*types.User, error) {
	_, err := s.db.Exec("UPDATE users SET first_name = ?, last_name = ?, email = ? WHERE id = ?", payload.FirstName, payload.LastName, payload.Email, id)
	if err != nil {
		return nil, err
	}

	user, err := s.GetUserByEmail(payload.Email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Store) UpdatePassword(id string, hashedPassword string) error {
	_, err := s.db.Exec("UPDATE users SET password = ? WHERE id = ?", hashedPassword, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) UpdateAvatar(id string, imageUrl string) (*types.User, error) {
	_, err := s.db.Exec("UPDATE users SET avatar = ? WHERE id = ?", imageUrl, id)
	if err != nil {
		return nil, err
	}

	user, err := s.GetUserById(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
