package request

type LikePostRequest struct {
	PostId int `json:"postId" binding:"required"`
}

type UnlikePostRequest struct {
	PostId int `json:"postId" binding:"required"`
}
