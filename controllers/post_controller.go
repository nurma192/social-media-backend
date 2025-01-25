package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"social-media-back/models/request"
	"social-media-back/models/response"
	"strconv"
)

func (c *AppController) CreatePost(ctx *gin.Context) {
	userId := ctx.MustGet("userId").(string)
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
	res, code, errRes := c.AppService.CreatePost(req, userId)
	if errRes != nil {
		ctx.IndentedJSON(code, errRes)
		return
	}

	ctx.IndentedJSON(code, res)
}

func (c *AppController) GetPost(ctx *gin.Context) {
	id := ctx.Param("id")
	userId := ctx.MustGet("userId").(string)
	postId, err := strconv.Atoi(id)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, response.DefaultResponse{
			Message: "Invalid post id",
			Detail:  err.Error(),
		})
		return
	}

	res, code, errRes := c.AppService.GetPostById(postId, userId)
	if errRes != nil {
		ctx.IndentedJSON(code, errRes)
		return
	}
	ctx.IndentedJSON(code, res)
}

func (c *AppController) GetAllPosts(ctx *gin.Context) {
	limitParam := ctx.DefaultQuery("limit", "10")
	pageParam := ctx.DefaultQuery("page", "1")
	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Limit must be a positive integer"})
		return
	}

	page, err := strconv.Atoi(pageParam)
	if err != nil || page <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Page must be a positive integer"})
		return
	}

	res, code, errRes := c.AppService.GetAllPosts(limit, page)

	if errRes != nil {
		ctx.IndentedJSON(code, errRes)
		return
	}
	ctx.IndentedJSON(code, res)
}

func (c *AppController) DeletePost(ctx *gin.Context) {
	id := ctx.Param("id")
	userId := ctx.MustGet("userId").(string)

	res, code, errRes := c.AppService.DeletePost(id, userId)
	if errRes != nil {
		ctx.IndentedJSON(code, errRes)
		return
	}

	ctx.IndentedJSON(code, res)
}
func (c *AppController) UpdatePost(ctx *gin.Context) {
	id := ctx.Param("id")
	userId := ctx.MustGet("userId").(string)

	postId, err := strconv.Atoi(id)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, response.DefaultResponse{
			Message: "Invalid post id",
			Detail:  err.Error(),
		})
		return
	}

	var req *request.UpdatePostRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, response.DefaultResponse{
			Message: "Invalid request body",
			Detail:  err.Error(),
		})
		return
	}

	res, code, errRes := c.AppService.UpdatePost(postId, userId, req)

	if errRes != nil {
		ctx.IndentedJSON(code, errRes)
		return
	}

	ctx.IndentedJSON(code, res)
}
