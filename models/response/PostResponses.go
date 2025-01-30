package response

import "social-media-back/models"

type CreatePostResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Post    *models.Post `json:"post"`
}

type GetPostByIdResponse struct {
	Success bool                    `json:"success"`
	Post    *models.PostWithAllInfo `json:"post"`
}

type UpdatePostResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Post    *models.Post `json:"post"`
}
