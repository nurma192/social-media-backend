package services

import (
	"net/http"
	"social-media-back/models/response"
)

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
