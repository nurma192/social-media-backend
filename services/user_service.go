package services

import (
	"net/http"
	"social-media-back/models/response"
)

func (s *AppService) CurrentUser(email string) (*response.CurrentUserResponse, int, *response.DefaultResponse) {
	user, err := s.DBService.GetUserByEmail(email)
	if err != nil {
		return nil, http.StatusInternalServerError, &response.DefaultResponse{
			Message: "Server Error",
			Detail:  "Error when try to get user",
		}
	}

	return &response.CurrentUserResponse{
		User:    user,
		Success: true,
	}, http.StatusOK, nil
}
