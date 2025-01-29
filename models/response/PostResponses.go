package response

import "social-media-back/models"

type CreatePostResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Post    *models.Post `json:"post"`
}

type GetPostsResponse struct {
	Success    bool                      `json:"success"`
	Page       int                       `json:"page"`
	TotalPages int                       `json:"totalPages"`
	Limit      int                       `json:"limit"`
	Posts      []*models.PostWithAllInfo `json:"posts"`
}

type UpdatePostResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Post    *models.Post `json:"post"`
}
