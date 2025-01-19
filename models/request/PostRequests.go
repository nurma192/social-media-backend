package request

import "mime/multipart"

type CreatePostRequest struct {
	ContentText string                  `form:"contentText"`
	Images      []*multipart.FileHeader `form:"images"`
}

type DeletePostRequest struct {
	PostId int `json:"postId" binding:"required"`
}

type UpdatePostRequest struct {
	ContentText   string                  `form:"contentText"`
	DeletedImages []string                `form:"deletedImages"`
	NewImages     []*multipart.FileHeader `form:"newImages"`
}
