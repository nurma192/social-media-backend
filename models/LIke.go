package models

type Like struct {
	Id        string `json:"id"`
	PostId    string `json:"post_id"`
	UserId    string `json:"user_id"`
	CreatedAt string `json:"created_at"`
}
