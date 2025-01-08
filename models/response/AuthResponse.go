package response

import "social-media-back/models"

type RegisterSuccessResponse struct {
	User    *models.User `json:"user"`
	Message string       `json:"message"`
	Success bool         `json:"success" `
}

type SendVerifyCodeResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type VerifyAccountResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}
