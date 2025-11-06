package models

import "time"

type Tweet struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CreateTweetRequest struct {
	UserID  int64  `json:"user_id"`
	Content string `json:"content" binding:"required,min=1,max=280"`
}
