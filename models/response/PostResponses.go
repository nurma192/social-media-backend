package response

import "social-media-back/models"

type CreatePostResponse struct {
	Post *models.Post `json:"post"`
}

type GetPostByIdResponse struct {
	Post *models.PostWithAllInfo `json:"post"`
}

type GetAllPostsResponse struct {
	Posts []*models.PostWithAllInfo `json:"posts"`
}

type UpdatePostResponse struct {
	Post *models.Post `json:"post"`
}
