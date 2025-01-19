package response

import "social-media-back/models"

type CreatePostResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Post    *models.Post `json:"post"`
}

type GetPostsResponse struct {
	Success bool           `json:"success"`
	Posts   []*models.Post `json:"posts"`
	Page    int            `json:"page"`
	Limit   int            `json:"limit"`
}
