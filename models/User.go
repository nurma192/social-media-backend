package models

import (
	"time"
)

type User struct {
	Id          int        `json:"id"`
	Email       string     `json:"email"`
	Username    string     `json:"username"`
	Firstname   string     `json:"firstname"`
	Lastname    string     `json:"lastname"`
	Password    string     `json:"password"`
	AvatarURL   *string    `json:"avatar_url"`
	DateOfBirth *time.Time `json:"date_of_birth"`
	Bio         *string    `json:"bio"`
	Verified    bool       `json:"verified"`
	Location    *string    `json:"location"`
	CreatedAt   time.Time  `json:"created_at"`
}
type UserFillInto struct {
	Id          int        `json:"id"`
	Email       string     `json:"email"`
	Username    string     `json:"username"`
	Firstname   string     `json:"firstname"`
	Lastname    string     `json:"lastname"`
	Password    string     `json:"password"`
	AvatarURL   *string    `json:"avatar_url"`
	DateOfBirth *time.Time `json:"date_of_birth"`
	Bio         *string    `json:"bio"`
	Verified    bool       `json:"verified"`
	Location    *string    `json:"location"`
	CreatedAt   time.Time  `json:"created_at"`
	PostsCount  int        `json:"posts_count"`
}
type UserMainInfo struct {
	Id        int     `json:"id"`
	Username  string  `json:"username"`
	Firstname string  `json:"firstname"`
	Lastname  string  `json:"lastname"`
	AvatarURL *string `json:"avatar_url"`
}
