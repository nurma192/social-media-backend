package models

type Comment struct {
	Id        int    `json:"id"`
	PostId    int    `json:"post_id"`
	UserId    string `json:"user_id"`
	Content   string `json:"text"`
	CreatedAt string `json:"created_at"`
}

type CommentWithUser struct {
	Id        int           `json:"id"`
	PostId    int           `json:"post_id"`
	User      *UserMainInfo `json:"user"`
	Content   string        `json:"text"`
	CreatedAt string        `json:"created_at"`
}
