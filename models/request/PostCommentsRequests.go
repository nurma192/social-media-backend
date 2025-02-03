package request

type CreateCommentRequest struct {
	Content string `json:"content"  binding:"required"`
	PostId  int    `json:"postId" binding:"required"`
}
type DeleteCommentRequest struct {
	CommentId int `json:"commentId" binding:"required"`
}
type UpdateCommentRequest struct {
	CommentId int    `json:"commentId" binding:"required"`
	Content   string `json:"content" binding:"required"`
}
type GetPostCommentsRequest struct {
	PostId int `json:"postId" binding:"required"`
}
