package response

import "social-media-back/models"

type LoginResponse struct {
	Token string `json:"token"`
}

type VerifyAccountResponse struct {
	User *models.User `json:"user"`
}
