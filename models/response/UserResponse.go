package response

import "social-media-back/models"

type CurrentUserResponse struct {
	User *models.User `json:"user"`
}
