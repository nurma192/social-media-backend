package models

import "time"

type Post struct {
	Id            string    `json:"id"`
	UserId        string    `json:"user_id"`
	ContentText   string    `json:"content_text"`
	Images        []Image   `json:"images"`
	LikesCount    int       `json:"likes_count"`
	CommentsCount int       `json:"comments_count"`
	CreatedAt     time.Time `json:"created_at"`
}

type PostWithAllInfo struct {
	Id            string        `json:"id"`
	User          *UserMainInfo `json:"user"`
	ContentText   string        `json:"content_text"`
	Images        []Image       `json:"images"`
	LikesCount    int           `json:"likes_count"`
	CommentsCount int           `json:"comments_count"`
	CreatedAt     time.Time     `json:"created_at"`
}
