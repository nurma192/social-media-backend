package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"social-media-back/models/request"
	"social-media-back/models/response"
)

func (c *AppController) LikePost(ctx *gin.Context) {
	userId := ctx.MustGet("userId").(string)
	var req = request.LikePostRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, response.DefaultResponse{
			Message: "invalid request body",
			Detail:  err.Error(),
		})
		return
	}

	res, code, errRes := c.AppService.LikePost(req.PostId, userId)
	if errRes != nil {
		ctx.IndentedJSON(code, errRes)
		return
	}

	ctx.IndentedJSON(code, res)
}

func (c *AppController) UnlikePost(ctx *gin.Context) {
	userId := ctx.MustGet("userId").(string)
	var req = request.UnlikePostRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, response.DefaultResponse{
			Message: "invalid request body",
			Detail:  err.Error(),
		})
		return
	}

	res, code, errRes := c.AppService.UnlikePost(req.PostId, userId)
	if errRes != nil {
		ctx.IndentedJSON(code, errRes)
		return
	}

	ctx.IndentedJSON(code, res)
}
