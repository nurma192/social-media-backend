package response

import "social-media-back/models"

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	Success      bool   `json:"success"`
}

type VerifyAccountResponse struct {
	User    *models.User `json:"user"`
	Message string       `json:"message"`
	Success bool         `json:"success" `
}
type VerifyAccountResponse2 struct {
	User    *models.User `json:"user"`
	Message string       `json:"message"`
	Success bool         `json:"success" `
}
