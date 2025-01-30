package models

type Comment struct {
	Id        int    `json:"id"`
	PostId    int    `json:"post_id"`
	UserId    int    `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type CommentWithUser struct {
	Id        int           `json:"id"`
	PostId    int           `json:"post_id"`
	User      *UserMainInfo `json:"user"`
	Content   string        `json:"content"`
	CreatedAt string        `json:"created_at"`
}
