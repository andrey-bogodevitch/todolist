package entity

import "time"

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	Login     string    `json:"login"`
	Password  string    `json:"password,omitempty"`
}
