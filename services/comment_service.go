package services

import (
	"net/http"
	"social-media-back/models/request"
	"social-media-back/models/response"
)

func (s *AppService) CreatePostComment(req *request.CreateCommentRequest, userId int) (*response.DefaultResponse, int, *response.DefaultResponse) {
	return nil, http.StatusOK, nil
}

func (s *AppService) DeletePostComment(req *request.DeleteCommentRequest) (*response.DefaultResponse, int, *response.DefaultResponse) {
	return nil, http.StatusOK, nil
}

func (s *AppService) UpdatePostComment(req *request.UpdateCommentRequest) (*response.DefaultResponse, int, *response.DefaultResponse) {
	return nil, http.StatusOK, nil
}

func (s *AppService) GetPostComments(postId, limit, page int) (*response.GetPostCommentsResponse, int, *response.DefaultResponse) {
	comments, err := s.DBService.GetPostComments(postId, limit, page)
	if err != nil {
		return nil, 0, &response.DefaultResponse{
			Message: err.Error(),
			Detail:  "services.GetPostComments",
		}
	}

	return &response.GetPostCommentsResponse{
		Comments: comments,
		Success:  true,
	}, http.StatusOK, nil
}
