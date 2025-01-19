package request

import "mime/multipart"

type CreatePostRequest struct {
	ContentText string                  `form:"contentText"`
	Images      []*multipart.FileHeader `form:"images"`
}

type DeletePostRequest struct {
	PostID int `json:"postId" binding:"required"`
}

type UpdatePostRequest struct {
	PostID        int      `form:"postId" binding:"required"`
	ContentText   string   `form:"contentText"`
	DeletedImages []string `form:"deletedImages"`
	NewImages     []string `form:"newImages"`
}
