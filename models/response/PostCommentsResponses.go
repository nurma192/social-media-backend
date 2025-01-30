package response

import "social-media-back/models"

type GetPostCommentsResponse struct {
	Comments   []models.CommentWithUser `json:"comments"`
	Page       int                      `json:"page"`
	TotalPages int                      `json:"totalPages"`
	Limit      int                      `json:"limit"`
	Success    bool                     `json:"success"`
}
