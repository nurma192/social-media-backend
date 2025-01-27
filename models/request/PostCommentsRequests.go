package request

type CreateCommentRequest struct {
	Content string `json:"content"`
	PostId  int    `json:"postId"`
}
type DeleteCommentRequest struct {
	CommentId int `json:"commentId"`
}
type UpdateCommentRequest struct {
	CommentId int    `json:"commentId"`
	Content   string `json:"content"`
}
type GetPostCommentsRequest struct {
	PostId int `json:"postId"`
}
