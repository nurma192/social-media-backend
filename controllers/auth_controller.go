package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"social-media-back/models/request"
	"social-media-back/models/response"
)

func (c *AppController) Login(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": c.AppService.Login("User", "password")})
}

func (c *AppController) Register(ctx *gin.Context) {
	var req request.RegisterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.DefaultErrorResponse{
			Detail:  err.Error(),
			Message: "Invalid request body",
		})
		return
	}

	message, user, status, err := c.AppService.Register(req)
	if err != nil {
		ctx.JSON(status, response.DefaultErrorResponse{
			Message: message,
			Detail:  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusBadRequest, response.RegisterSuccessResponse{
		User:    user,
		Success: true,
		Message: message,
	})

}

func (c *AppController) SendVerifyCode(ctx *gin.Context) {
	var req request.SendVerifyCodeRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.DefaultErrorResponse{
			Detail:  err.Error(),
			Message: "Invalid request body",
		})
		return
	}

	message, code, err := c.AppService.SendVerifyCode(req.Email)
	if err != nil {
		ctx.JSON(code, response.DefaultErrorResponse{
			Message: message,
			Detail:  err.Error(),
		})
		return
	}

	ctx.JSON(code, response.SendVerifyCodeResponse{
		Success: true,
		Message: message,
	})
}

func (c *AppController) VerifyAccount(ctx *gin.Context) {
	var req request.VerifyAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.DefaultErrorResponse{
			Detail:  err.Error(),
			Message: "Invalid request body",
		})
		return
	}

	message, code, err := c.AppService.VerifyAccount(req.Email, req.Code)
	if err != nil {
		ctx.JSON(code, response.DefaultErrorResponse{
			Message: message,
			Detail:  err.Error(),
		})
		return
	}

	ctx.JSON(code, response.DefaultSuccessResponse{
		Success: true,
		Message: message,
	})
}
