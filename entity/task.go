package entity

import "time"

type Task struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Task      string    `json:"task"`
	Status    string    `json:"status"`
	DeletedAt time.Time `json:"deleted_at"`
}
