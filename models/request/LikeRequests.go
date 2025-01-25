package request

type LikePostRequest struct {
	PostId int `json:"post_id" binding:"required"`
}

type UnlikePostRequest struct {
	PostId int `json:"post_id" binding:"required"`
}
