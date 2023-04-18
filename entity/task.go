package entity

type Task struct {
	UserID int64  `json:"user_id"`
	Task   string `json:"task"`
	Status string `json:"status"`
}
