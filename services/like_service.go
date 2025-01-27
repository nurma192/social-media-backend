package services

import (
	"net/http"
	"social-media-back/models/response"
)

func (s *AppService) LikePost(postId, userId int) (*response.DefaultResponse, int, *response.DefaultResponse) {
	err := s.DBService.AddLikePost(postId, userId)
	if err != nil {
		return nil, http.StatusInternalServerError, &response.DefaultResponse{
			Message: "Failed to add like post",
			Detail:  err.Error(),
		}
	}
	return &response.DefaultResponse{
		Success: true,
	}, http.StatusCreated, nil
}

func (s *AppService) UnlikePost(postId, userId int) (*response.DefaultResponse, int, *response.DefaultResponse) {
	err := s.DBService.DeleteLikePost(postId, userId)
	if err != nil {
		return nil, http.StatusInternalServerError, &response.DefaultResponse{
			Message: "Failed to delete like post",
			Detail:  err.Error(),
		}
	}
	return &response.DefaultResponse{
		Success: true,
	}, http.StatusCreated, nil
}
