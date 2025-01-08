package response

import "social-media-back/models"

type RegisterSuccessResponse struct {
	User    *models.User `json:"user"`
	Message string       `json:"message"`
	Success bool         `json:"success" `
}

type RegisterErrorResponse struct {
	Message string `json:"message"`
	Detail  string `json:"detail"`
	Success bool   `json:"success" `
}
