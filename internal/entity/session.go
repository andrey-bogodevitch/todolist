package entity

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        uuid.UUID `json:"id"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiredAt time.Time `json:"expired_at"`
}
