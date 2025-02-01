package services

import (
	"net/http"
	"social-media-back/models/response"
)

func (s *AppService) CurrentUser(email string) (*response.Response, int) {
	user, err := s.DBService.GetUserByEmail(email)
	if err != nil {
		return &response.Response{
			Error: err.Error(),
		}, http.StatusInternalServerError
	}

	currentUserRes := &response.CurrentUserResponse{
		User: user,
	}

	return &response.Response{
		Result: currentUserRes,
	}, http.StatusOK
}
