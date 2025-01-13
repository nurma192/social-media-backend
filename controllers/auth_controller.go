package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"social-media-back/models/request"
	"social-media-back/models/response"
)

func (c *AppController) Login(ctx *gin.Context) {
	var req request.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.DefaultErrorResponse{
			Message: "Bad request",
			Detail:  err.Error(),
		})
		return
	}

	res, code, errRes := c.AppService.Login(req.Email, req.Password)
	if errRes != nil {
		ctx.JSON(code, errRes)
		return
	}

	ctx.SetCookie("RefreshToken", res.RefreshToken, 3600*24*7, "", "", false, true)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"token":   res.Token,
	})
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

	res, status, errRes := c.AppService.Register(req)
	if errRes != nil {
		ctx.JSON(status, errRes)
		return
	}

	ctx.JSON(http.StatusCreated, res)
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

	res, code, errRes := c.AppService.SendVerifyCode(req.Email)
	if errRes != nil {
		ctx.JSON(code, errRes)
		return
	}

	ctx.JSON(code, res)
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

	res, code, errRes := c.AppService.VerifyAccount(req.Email, req.Code)
	if errRes != nil {
		ctx.JSON(code, errRes)
		return
	}

	ctx.JSON(code, res)
}

func (c *AppController) RefreshToken(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("RefreshToken")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.DefaultErrorResponse{
			Message: "Refresh token not found in cookie",
			Detail:  err.Error(),
		})
		return
	}

	res, code, errRes := c.AppService.RefreshToken(refreshToken)
	if errRes != nil {
		ctx.JSON(code, errRes)
		return
	}

	ctx.SetCookie("RefreshToken", res.RefreshToken, 3600*24*7, "", "", false, true)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"token":   res.Token,
	})

}
