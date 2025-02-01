package services

import (
	"fmt"
	"net/http"
	"social-media-back/models/request"
	"social-media-back/models/response"
)

func (s *AppService) CreatePostComment(req *request.CreateCommentRequest, userId int) (*response.Response, int) {
	err := s.DBService.CreatePostComment(req.Content, req.PostId, userId)
	if err != nil {
		return &response.Response{
			Error: err.Error(),
		}, http.StatusInternalServerError
	}
	return &response.Response{}, http.StatusCreated
}

func (s *AppService) DeletePostComment(req *request.DeleteCommentRequest, userId int) (*response.Response, int) {
	err := s.DBService.DeletePostComment(req.CommentId, userId)
	if err != nil {
		return &response.Response{
			Error: err.Error(),
		}, http.StatusInternalServerError
	}
	return &response.Response{}, http.StatusOK
}

func (s *AppService) UpdatePostComment(req *request.UpdateCommentRequest, userId int) (*response.Response, int) {
	err := s.DBService.UpdatePostComment(req.CommentId, req.Content, userId)
	if err != nil {
		return &response.Response{
			Error: err.Error(),
		}, http.StatusInternalServerError
	}
	return &response.Response{}, http.StatusOK
}

func (s *AppService) GetPostComments(postId, limit, page int) (*response.Response, int) {
	comments, totalPages, err := s.DBService.GetPostComments(postId, limit, page)
	if err != nil {
		return &response.Response{
			Error: err.Error(),
		}, http.StatusInternalServerError
	}

	fmt.Println(comments)

	getPostCommentsRes := &response.GetPostCommentsResponse{
		Comments: comments,
	}

	paginationRes := response.PaginationResponse{
		Result:     getPostCommentsRes,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}

	return &response.Response{
		Result: paginationRes,
	}, http.StatusOK
}
