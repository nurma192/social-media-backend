package request

import "mime/multipart"

type CreatePostRequest struct {
	UserID      string                  `form:"userId" binding:"required"`
	ContentText string                  `form:"contentText"`
	Images      []*multipart.FileHeader `form:"images"`
}
