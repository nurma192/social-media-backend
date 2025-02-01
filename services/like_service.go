package services

import (
	"net/http"
	"social-media-back/models/response"
)

func (s *AppService) LikePost(postId, userId int) (*response.Response, int) {
	err := s.DBService.AddLikePost(postId, userId)
	if err != nil {
		return &response.Response{
			Error: err.Error(),
		}, http.StatusInternalServerError
	}
	return &response.Response{}, http.StatusCreated
}

func (s *AppService) UnlikePost(postId, userId int) (*response.Response, int) {
	err := s.DBService.DeleteLikePost(postId, userId)
	if err != nil {
		return &response.Response{
			Error: err.Error(),
		}, http.StatusInternalServerError
	}
	return &response.Response{}, http.StatusCreated
}
