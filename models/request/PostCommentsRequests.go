package request

type CreateCommentRequest struct {
	Content string `json:"content"`
	PostId  string `json:"postId"`
}
type DeleteCommentRequest struct {
	CommentId string `json:"commentId"`
}
type UpdateCommentRequest struct {
	CommentId string `json:"commentId"`
	Content   string `json:"content"`
}
type GetPostCommentsRequest struct {
	PostId string `json:"postId"`
}
