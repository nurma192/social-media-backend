package response

import "social-media-back/models"

type CurrentUserResponse struct {
	Message string       `json:"message"`
	User    *models.User `json:"user"`
	Success bool         `json:"success"`
}
