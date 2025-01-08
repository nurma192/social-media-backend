package response

import "social-media-back/models"

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	Success      bool   `json:"success"`
}

type RegisterResponse struct {
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
