package response

import "social-media-back/models"

type CreatePostResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Post    *models.Post `json:"post"`
}
