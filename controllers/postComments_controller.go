package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"social-media-back/models/request"
	"social-media-back/models/response"
	"strconv"
)

func (c *AppController) CreatePostComment(ctx *gin.Context) {
	userId := ctx.MustGet("userId").(int)
	var req request.CreateCommentRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, response.DefaultResponse{
			Message: "invalid request body",
			Detail:  err.Error(),
		})
	}

	res, code := c.AppService.CreatePostComment(&req, userId)

	ctx.IndentedJSON(code, res)

}
func (c *AppController) DeletePostComment(ctx *gin.Context) {
	userId := ctx.MustGet("userId").(int)
	var req request.DeleteCommentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, response.DefaultResponse{
			Message: "invalid request body",
			Detail:  err.Error(),
		})
	}

	res, code := c.AppService.DeletePostComment(&req, userId)

	ctx.IndentedJSON(code, res)

}
func (c *AppController) UpdatePostComment(ctx *gin.Context) {
	userId := ctx.MustGet("userId").(int)

	var req request.UpdateCommentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, response.DefaultResponse{
			Message: "invalid request body",
			Detail:  err.Error(),
		})
		return
	}

	res, code := c.AppService.UpdatePostComment(&req, userId)

	ctx.IndentedJSON(code, res)

}
func (c *AppController) GetPostsComments(ctx *gin.Context) {
	id := ctx.Param("id")
	postId, err := strconv.Atoi(id)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, response.DefaultResponse{
			Message: "Invalid post id",
			Detail:  err.Error(),
		})
		return
	}
	limitParam := ctx.DefaultQuery("limit", "10")
	pageParam := ctx.DefaultQuery("page", "1")
	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Limit must be a positive integer"})
		return
	}

	page, err := strconv.Atoi(pageParam)
	if err != nil || page <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Page must be a positive integer"})
		return
	}

	res, code := c.AppService.GetPostComments(postId, limit, page)

	ctx.IndentedJSON(code, res)
}
