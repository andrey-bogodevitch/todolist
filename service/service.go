package service

import (
	"todolist/storage"

	"golang.org/x/crypto/bcrypt"
)

type Storage interface {
	CreateUser(name, login, password string) error
	GetUserByID(id int64) (storage.User, error)
}

type User struct {
	storage Storage
}

func NewUser(s Storage) *User {
	return &User{storage: s}
}

func (u *User) AddUser(name, login, password string) error {
	passwordByte := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(passwordByte, 10)
	if err != nil {
		return err
	}

	err = u.storage.CreateUser(name, login, string(hash))
	if err != nil {
		return err
	}
	return nil
}

func (u *User) GetUser(id int64) (storage.User, error) {
	user, err := u.storage.GetUserByID(id)
	if err != nil {
		return storage.User{}, err
	}
	return user, nil
}
