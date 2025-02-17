package models

import "time"

type Comment struct {
	Id        int    `json:"id"`
	PostId    int    `json:"postId"`
	UserId    int    `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type CommentWithUser struct {
	Id        int           `json:"id"`
	PostId    int           `json:"postId"`
	User      *UserMainInfo `json:"user"`
	Content   string        `json:"content"`
	CreatedAt time.Time     `json:"createdAt"`
}
