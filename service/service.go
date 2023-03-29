package service

import (
	"time"

	"todolist/entity"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Storage interface {
	CreateUser(user entity.User) error
	GetUserByID(id int64) (entity.User, error)
	UserByLogin(login string) (entity.User, error)
	SaveSession(session entity.Session) error
	SessionByID(id uuid.UUID) (entity.Session, error)
	DeleteUser(id int64) error
}

type User struct {
	storage Storage
}

func NewUser(s Storage) *User {
	return &User{storage: s}
}

func (u *User) AddUser(user entity.User) error {
	passwordByte := []byte(user.Password)
	hash, err := bcrypt.GenerateFromPassword(passwordByte, 10)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	user.CreatedAt = time.Now()
	user.Role = "user"

	err = u.storage.CreateUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) GetUser(id int64) (entity.User, error) {
	user, err := u.storage.GetUserByID(id)
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (u *User) CreateSession(login, password string) (entity.Session, error) {
	user, err := u.storage.UserByLogin(login)
	if err != nil {
		return entity.Session{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return entity.Session{}, err
	}

	sessionID, err := uuid.NewUUID()
	if err != nil {
		return entity.Session{}, err
	}

	now := time.Now()
	session := entity.Session{
		ID:        sessionID,
		UserID:    user.ID,
		CreatedAt: now,
		ExpiredAt: now.Add(10 * time.Minute),
	}

	err = u.storage.SaveSession(session)
	if err != nil {
		return entity.Session{}, err
	}

	return session, nil
}

func (u *User) FindSessionByID(id uuid.UUID) (entity.Session, error) {
	session, err := u.storage.SessionByID(id)
	if err != nil {
		return entity.Session{}, err
	}

	return session, nil
}

func (u *User) DeleteUser(id int64) error {
	err := u.storage.DeleteUser(id)
	if err != nil {
		return err
	}

	return nil
}
