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
		ctx.IndentedJSON(http.StatusBadRequest, response.DefaultResponse{
			Message: "Bad request",
			Detail:  err.Error(),
		})
		return
	}

	res, code, refreshToken := c.AppService.Login(req.Email, req.Password)

	ctx.SetCookie("RefreshToken", refreshToken, 3600*24*7, "", "", false, true)
	ctx.IndentedJSON(code, res)
}

func (c *AppController) Register(ctx *gin.Context) {
	var req request.RegisterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, response.DefaultResponse{
			Detail:  err.Error(),
			Message: "Invalid request body",
		})
		return
	}

	res, status := c.AppService.Register(req)

	ctx.IndentedJSON(status, res)
}

func (c *AppController) SendVerifyCode(ctx *gin.Context) {
	var req request.SendVerifyCodeRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, response.DefaultResponse{
			Detail:  err.Error(),
			Message: "Invalid request body",
		})
		return
	}

	res, code := c.AppService.SendVerifyCode(req.Email)

	ctx.IndentedJSON(code, res)
}

func (c *AppController) VerifyAccount(ctx *gin.Context) {
	var req request.VerifyAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, response.DefaultResponse{
			Detail:  err.Error(),
			Message: "Invalid request body",
		})
		return
	}

	res, code := c.AppService.VerifyAccount(req.Email, req.Code)

	ctx.IndentedJSON(code, res)
}

func (c *AppController) RefreshToken(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("RefreshToken")
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, response.DefaultResponse{
			Message: "Refresh token not found in cookie",
			Detail:  err.Error(),
		})
		return
	}

	res, code, newRefreshToken := c.AppService.RefreshToken(refreshToken)

	ctx.SetCookie("RefreshToken", newRefreshToken, 3600*24*7, "", "", false, true)
	ctx.IndentedJSON(code, res)

}
