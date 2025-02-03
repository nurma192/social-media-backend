package response

import "social-media-back/models"

type GetPostCommentsResponse struct {
	Comments []models.CommentWithUser `json:"comments"`
}

type CreatePostCommentResponse struct {
	Comment models.CommentWithUser `json:"comment"`
}
