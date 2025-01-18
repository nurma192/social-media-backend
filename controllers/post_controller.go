package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/http"
	"social-media-back/models/response"
)

type CreatePostRequest struct {
	UserID      string                  `form:"userId" binding:"required"`
	ContentText string                  `form:"contentText"`
	Images      []*multipart.FileHeader `form:"images"`
}

func (c *AppController) CreatePost(ctx *gin.Context) {
	var req CreatePostRequest

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, response.DefaultResponse{
			Message: "Invalid request body",
			Detail:  err.Error(),
		})
		return
	}
	if len(req.Images) > 5 {
		ctx.IndentedJSON(http.StatusBadRequest, response.DefaultResponse{
			Message: "Max Image size is 5",
		})
		return
	}
	fmt.Println(req)

	var uploadedImagesURLs []string

	for _, fileHeader := range req.Images {
		file, err := fileHeader.Open()
		if err != nil {
			ctx.IndentedJSON(http.StatusBadRequest, response.DefaultResponse{
				Message: "Failed to open image",
				Detail:  err.Error(),
			})
			return
		}
		defer file.Close()
		fmt.Println(fileHeader.Filename)

		fileURL, err := c.AppService.AWSService.UploadFile(file, fileHeader.Filename, fileHeader.Header.Get("Content-Type"))
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, response.DefaultResponse{
				Message: "Failed to upload image to S3",
				Detail:  err.Error(),
			})
			return
		}
		uploadedImagesURLs = append(uploadedImagesURLs, fileURL)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Post created successfully",
		"userId":  req.UserID,
		"images":  uploadedImagesURLs, // Возвращаем загруженные URL
	})
}

func (c *AppController) GetPost(ctx *gin.Context) {
	ctx.IndentedJSON(200, response.DefaultResponse{
		Message: "Get Post",
	})
}

func (c *AppController) GetAllPosts(ctx *gin.Context) {
	ctx.IndentedJSON(200, response.DefaultResponse{
		Message: "Get All Posts",
	})
}

func (c *AppController) DeletePost(ctx *gin.Context) {
	ctx.IndentedJSON(200, response.DefaultResponse{
		Message: "Delete Post",
	})
}
func (c *AppController) UpdatePost(ctx *gin.Context) {
	ctx.IndentedJSON(200, response.DefaultResponse{
		Message: "Update Post",
	})
}
