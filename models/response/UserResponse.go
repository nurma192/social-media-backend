package response

import "social-media-back/models"

type GetCurrentUserResponse struct {
	User *models.User `json:"user"`
}

type GetUserResponse struct {
	User *models.User `json:"user"`
}
