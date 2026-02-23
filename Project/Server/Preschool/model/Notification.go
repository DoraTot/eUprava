package model

type Notification struct {
	ID        int    `json:"id" db:"id"`
	Message   string `json:"message" db:"message"`
	UserId    string `json:"user_id" db:"user_id"`
	IsRead    bool   `json:"read" db:"read"`
	CreatedAt string `json:"created_at" db:"created_at"`
}
