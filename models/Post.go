package models

import "time"

type Post struct {
	Id            int       `json:"id"`
	UserId        int       `json:"userId"`
	ContentText   string    `json:"contentText"`
	Images        []Image   `json:"images"`
	LikesCount    int       `json:"likesCount"`
	CommentsCount int       `json:"commentsCount"`
	CreatedAt     time.Time `json:"created_at"`
}

type PostWithAllInfo struct {
	Id            int           `json:"id"`
	User          *UserMainInfo `json:"user"`
	ContentText   string        `json:"contentText"`
	Images        []Image       `json:"images"`
	LikedByUser   bool          `json:"likedByUser"`
	LikesCount    int           `json:"likesCount"`
	CommentsCount int           `json:"commentsCount"`
	CreatedAt     time.Time     `json:"createdAt"`
}
