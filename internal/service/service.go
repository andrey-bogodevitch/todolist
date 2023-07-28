package service

import (
	"context"
	"time"
	entity2 "todolist/internal/entity"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	StatusActive    = "active"
	StatusCompleted = "completed"
	StatusDeleted   = "deleted"
)

type Storage interface {
	CreateUser(user entity2.User) error
	GetUserByID(id int64) (entity2.User, error)
	UserByLogin(login string) (entity2.User, error)
	SaveSession(session entity2.Session) error
	SessionByID(id uuid.UUID) (entity2.Session, error)
	DeleteUser(id int64) error
	AddAdminRules(id int64) error
	CreateTask(task entity2.Task) error
	UpdateTask(task entity2.Task, status string) error
	GetTasksByUserID(id int64) ([]entity2.Task, error)
	GetTaskByID(id int64) (entity2.Task, error)
	DeleteTask(taskID int64, status string) error
	SaveSessionRedis(ctx context.Context, session entity2.Session) error
	SessionByIDRedis(ctx context.Context, id uuid.UUID) (entity2.Session, error)
}

type User struct {
	storage Storage
}

func NewUser(s Storage) *User {
	return &User{storage: s}
}

func (u *User) AddUser(user entity2.User) error {
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

func (u *User) GetUser(id int64) (entity2.User, error) {
	user, err := u.storage.GetUserByID(id)
	if err != nil {
		return entity2.User{}, err
	}
	return user, nil
}

func (u *User) CreateSession(ctx context.Context, login, password string) (entity2.Session, error) {
	user, err := u.storage.UserByLogin(login)
	if err != nil {
		return entity2.Session{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return entity2.Session{}, err
	}

	sessionID, err := uuid.NewUUID()
	if err != nil {
		return entity2.Session{}, err
	}

	now := time.Now()
	session := entity2.Session{
		ID:        sessionID,
		UserID:    user.ID,
		CreatedAt: now,
		ExpiredAt: now.Add(10 * time.Minute),
	}

	err = u.storage.SaveSessionRedis(ctx, session)
	if err != nil {
		return entity2.Session{}, err
	}

	return session, nil
}

func (u *User) FindSessionByID(ctx context.Context, id uuid.UUID) (entity2.Session, error) {
	session, err := u.storage.SessionByIDRedis(ctx, id)
	if err != nil {
		return entity2.Session{}, err
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

func (u *User) AddAdminRules(id int64) error {
	err := u.storage.AddAdminRules(id)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) AddTask(task entity2.Task) error {
	task.Status = StatusActive
	err := u.storage.CreateTask(task)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) UpdateTask(task entity2.Task) error {
	err := u.storage.UpdateTask(task, StatusCompleted)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) GetTasks(id int64) ([]entity2.Task, error) {
	tasks, err := u.storage.GetTasksByUserID(id)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (u *User) GetTaskByID(id int64) (entity2.Task, error) {
	task, err := u.storage.GetTaskByID(id)
	if err != nil {
		return entity2.Task{}, err
	}
	return task, nil
}

func (u *User) DeleteTask(taskID int64) error {
	err := u.storage.DeleteTask(taskID, StatusDeleted)
	if err != nil {
		return err
	}
	return nil
}
