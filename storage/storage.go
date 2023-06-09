package storage

import (
	"database/sql"
	"time"

	"todolist/entity"

	"github.com/google/uuid"
)

type UserStorage struct {
	db *sql.DB
}

func NewUserStorage(dbpool *sql.DB) *UserStorage {
	return &UserStorage{
		db: dbpool,
	}
}

func (s *UserStorage) CreateUser(user entity.User) error {
	query := "INSERT INTO users (name, role, created_at, login, password) values ($1, $2, $3, $4, $5)"
	_, err := s.db.Exec(query, user.Name, user.Role, user.CreatedAt, user.Login, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserStorage) GetUserByID(id int64) (entity.User, error) {
	query := "SELECT id, name, role, created_at, login FROM users WHERE id = $1 AND deleted_at IS NULL"

	var user entity.User

	err := s.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Role,
		&user.CreatedAt,
		&user.Login,
	)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}
func (s *UserStorage) UserByLogin(login string) (entity.User, error) {
	query := "SELECT id, name, role, created_at, login, password FROM users WHERE login = $1 AND deleted_at IS NULL"

	var user entity.User

	err := s.db.QueryRow(query, login).Scan(
		&user.ID,
		&user.Name,
		&user.Role,
		&user.CreatedAt,
		&user.Login,
		&user.Password,
	)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (s *UserStorage) SaveSession(session entity.Session) error {
	query := "INSERT INTO sessions (id, user_id, created_at, expired_at) values ($1, $2, $3, $4)"
	_, err := s.db.Exec(query, session.ID, session.UserID, session.CreatedAt, session.ExpiredAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserStorage) SessionByID(id uuid.UUID) (entity.Session, error) {
	var session entity.Session
	query := "SELECT id, user_id, created_at, expired_at FROM sessions WHERE id = $1"
	err := s.db.QueryRow(query, id).Scan(
		&session.ID,
		&session.UserID,
		&session.CreatedAt,
		&session.ExpiredAt,
	)
	if err != nil {
		return entity.Session{}, err
	}

	return session, nil
}

func (s *UserStorage) DeleteUser(id int64) error {
	query := "UPDATE users SET deleted_at = $1 where id = $2 AND deleted_at IS NULL"
	_, err := s.db.Exec(query, time.Now(), id)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserStorage) AddAdminRules(id int64) error {
	query := "UPDATE users SET role = $1 where id = $2 AND deleted_at IS NULL"
	_, err := s.db.Exec(query, "admin", id)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserStorage) CreateTask(task entity.Task) error {
	query := "INSERT INTO tasks (user_id, task, status) values ($1, $2, $3)"
	_, err := s.db.Exec(query, task.UserID, task.Task, task.Status)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserStorage) UpdateTask(task entity.Task, status string) error {
	query := "UPDATE tasks SET status = $1 WHERE  id = $2"
	_, err := s.db.Exec(query, status, task.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserStorage) GetTasksByUserID(id int64) ([]entity.Task, error) {
	query := "SELECT id, user_id, task, status FROM tasks WHERE user_id = $1 AND deleted_at IS NULL"

	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []entity.Task

	for rows.Next() {
		var task entity.Task
		err = rows.Scan(&task.ID, &task.UserID, &task.Task, &task.Status)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (s *UserStorage) GetTaskByID(id int64) (entity.Task, error) {
	query := "SELECT id, user_id, task, status FROM tasks WHERE id = $1 AND deleted_at IS NULL"

	var task entity.Task

	err := s.db.QueryRow(query, id).Scan(
		&task.ID,
		&task.UserID,
		&task.Task,
		&task.Status,
	)
	if err != nil {
		return entity.Task{}, err
	}

	return task, nil
}

func (s *UserStorage) DeleteTask(taskID int64, status string) error {
	query := "UPDATE tasks SET deleted_at = $1, status = $2 WHERE  id = $3"
	_, err := s.db.Exec(query, time.Now(), status, taskID)
	if err != nil {
		return err
	}
	return nil
}
