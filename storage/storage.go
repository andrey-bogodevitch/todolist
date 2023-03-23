package storage

import (
	"database/sql"
	"time"
)

type User struct {
	ID         int64
	Name       string
	Role       string
	Created_at time.Time
	Login      string
	Password   string
}

type UserStorage struct {
	db *sql.DB
}

func NewUserStorage(dbpool *sql.DB) *UserStorage {
	return &UserStorage{
		db: dbpool,
	}
}

func (s *UserStorage) CreateUser(name, login, password string) error {
	query := "INSERT INTO users (name, role, created_at, login, password) values ($1, $2, $3, $4, $5)"
	_, err := s.db.Exec(query, name, "user", time.Now(), login, password)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserStorage) GetUserByID(id int64) (User, error) {
	query := "SELECT id, name, role, created_at, login, password FROM users WHERE id = $1"

	var user User

	err := s.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Role,
		&user.Created_at,
		&user.Login,
		&user.Password,
	)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
