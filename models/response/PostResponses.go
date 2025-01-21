package response

import "social-media-back/models"

type CreatePostResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Post    *models.Post `json:"post"`
}

type GetPostsResponse struct {
	Success bool                   `json:"success"`
	Page    int                    `json:"page"`
	Limit   int                    `json:"limit"`
	Posts   []*models.PostWithUser `json:"posts"`
}

type UpdatePostResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Post    *models.Post `json:"post"`
}
