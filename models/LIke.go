package models

type Like struct {
	Id        int    `json:"id"`
	PostId    int    `json:"post_id"`
	UserId    int    `json:"user_id"`
	CreatedAt string `json:"created_at"`
}
