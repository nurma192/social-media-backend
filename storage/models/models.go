package models

import "time"

type User struct {
	ID          uint      `json:"id" gorm:"primary_key"`
	Email       string    `json:"email" gorm:"unique;not null"`
	Password    string    `json:"password" gorm:"not null"`
	Name        string    `json:"name"`
	AvatarURL   string    `json:"avatar_url"`
	DateOfBirth time.Time `json:"date_of_birth" gorm:"default:NULL"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Bio         string    `json:"bio"`
	Location    string    `json:"location"`
	Posts       []Post    `json:"posts"`
	Likes       []Like    `json:"likes"`
	Comments    []Comment `json:"comments"`
	Followers   []Follow  `json:"followers" gorm:"foreignKey:FollowingID;references:ID"`
	Following   []Follow  `json:"following" gorm:"foreignKey:FollowerID;references:ID"`
}

type Follow struct {
	ID          uint `json:"id" gorm:"primary_key"`
	FollowerID  uint `json:"follower_id" gorm:"not null"`
	FollowingID uint `json:"following_id" gorm:"not null"`
	Follower    User `json:"follower" gorm:"foreignKey:FollowerID;references:ID"`
	Following   User `json:"following" gorm:"foreignKey:FollowingID;references:ID"`
}

type Post struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	Content   string    `json:"content" gorm:"not null"`
	Author    User      `json:"author" gorm:"foreignKey:AuthorID;references:ID"`
	AuthorID  uint      `json:"author_id" gorm:"not null"`
	Likes     []Like    `json:"likes"`
	Comments  []Comment `json:"comments"`
	CreatedAt time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Like struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	User      User      `json:"user" gorm:"foreignKey:UserID;references:ID"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	Post      Post      `json:"post" gorm:"foreignKey:PostID;references:ID"`
	PostID    uint      `json:"post_id" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
}

type Comment struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	Content   string    `json:"content" gorm:"not null"`
	User      User      `json:"user" gorm:"foreignKey:UserID;references:ID"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	Post      Post      `json:"post" gorm:"foreignKey:PostID;references:ID"`
	PostID    uint      `json:"post_id" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
