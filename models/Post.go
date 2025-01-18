package models

import "time"

type Post struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	ContentText string    `json:"content_text"`
	Images      []string  `json:"images"`
	CreatedAt   time.Time `json:"created_at"`
}
