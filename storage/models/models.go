package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"CreatedAt"`
	UpdatedAt time.Time      `json:"UpdatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Email       string    `json:"email" gorm:"unique;not null"`
	Password    string    `json:"password" gorm:"not null"`
	Name        string    `json:"name"`
	AvatarURL   string    `json:"avatar_url"`
	DateOfBirth time.Time `json:"date_of_birth" gorm:"default:NULL"`
	Bio         string    `json:"bio"`
	Location    string    `json:"location"`
	Likes       []Like    `json:"likes"      `
	Posts       []Post    `json:"posts"      `
	Comments    []Comment `json:"comments"   `
	Followers   []User    `json:"followers"  gorm:"many2many:followers"`
	Followings  []User    `json:"followings" gorm:"many2many:followings"`
}

type Like struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"CreatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	UserID    uint           `json:"userID" gorm:"not null;uniqueIndex:idx_user_post"`
	PostID    uint           `json:"postID" gorm:"not null;uniqueIndex:idx_user_post"`
}

type Comment struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"CreatedAt"`
	UpdatedAt time.Time      `json:"UpdatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Content   string         `json:"content" gorm:"not null"`
	UserID    uint           `json:"userID" gorm:"not null"`
	PostID    uint           `json:"postID" gorm:"not null"`
}

type Post struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"CreatedAt"`
	UpdatedAt time.Time      `json:"UpdatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Content   string         `json:"content"  gorm:"size:255;not null"`
	UserID    uint           `json:"authorID" gorm:"not null"`
	Likes     []Like         `json:"likes"`
	Comments  []Comment      `json:"comments"`
}
