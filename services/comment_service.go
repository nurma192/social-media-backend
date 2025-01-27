package services

import (
	"net/http"
	"social-media-back/models/request"
	"social-media-back/models/response"
)

func (s *AppService) CreatePostComment(req *request.CreateCommentRequest, userId int) (*response.DefaultResponse, int, *response.DefaultResponse) {
	err := s.DBService.CreatePostComment(req.Content, req.PostId, userId)
	if err != nil {
		return nil, http.StatusInternalServerError, &response.DefaultResponse{
			Message: err.Error(),
			Detail:  "services.comment_service.CreatePostComment(1)",
		}
	}
	return nil, http.StatusCreated, &response.DefaultResponse{
		Success: true,
	}
}

func (s *AppService) DeletePostComment(req *request.DeleteCommentRequest, userId int) (*response.DefaultResponse, int, *response.DefaultResponse) {
	err := s.DBService.DeletePostComment(req.CommentId, userId)
	if err != nil {
		return nil, http.StatusInternalServerError, &response.DefaultResponse{
			Message: err.Error(),
			Detail:  "services.comment_service.DeletePostComment(1)",
		}
	}
	return nil, http.StatusOK, &response.DefaultResponse{
		Success: true,
	}
}

func (s *AppService) UpdatePostComment(req *request.UpdateCommentRequest, userId int) (*response.DefaultResponse, int, *response.DefaultResponse) {
	err := s.DBService.UpdatePostComment(req.CommentId, req.Content, userId)
	if err != nil {
		return nil, http.StatusInternalServerError, &response.DefaultResponse{
			Message: err.Error(),
			Detail:  "services.comment_service.UpdatePostComment(1)",
		}
	}
	return &response.DefaultResponse{
		Success: true,
	}, http.StatusOK, nil
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
