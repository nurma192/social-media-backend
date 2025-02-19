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

	currentUserRes := &response.GetCurrentUserResponse{
		User: user,
	}

	return &response.Response{
		Result: currentUserRes,
	}, http.StatusOK
}

func (s *AppService) GetUserById(id int) *response.Response {
	user, err := s.DBService.GetUserById(id)
	if err != nil {
		return &response.Response{
			Error: err.Error(),
			Code:  http.StatusInternalServerError,
		}
	}
	if user == nil {
		return &response.Response{
			Error: "user not found",
			Code:  http.StatusNotFound,
		}
	}

	getUserRes := &response.GetUserResponse{
		User: user,
	}

	return &response.Response{
		Result: getUserRes,
		Code:   http.StatusOK,
	}
}
