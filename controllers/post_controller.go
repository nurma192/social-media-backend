package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"social-media-back/models/request"
	"social-media-back/models/response"
)

func (c *AppController) CreatePost(ctx *gin.Context) {
	var req request.CreatePostRequest

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
	res, code, errRes := c.AppService.CreatePost(req)
	if errRes != nil {
		ctx.IndentedJSON(code, errRes)
		return
	}

	ctx.IndentedJSON(code, res)
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
