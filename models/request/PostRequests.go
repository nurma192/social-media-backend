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
	ContentText     string                  `form:"contentText"`
	DeletedImagesId []string                `form:"deletedImagesId"`
	NewImages       []*multipart.FileHeader `form:"newImages"`
}
