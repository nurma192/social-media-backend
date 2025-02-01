package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"social-media-back/models/response"
)

func (c *AppController) Current(ctx *gin.Context) {
	email := ctx.MustGet("email").(string)

	if email == "" {
		ctx.JSON(http.StatusBadRequest, response.DefaultResponse{
			Message: "Email not found",
			Detail:  "Email not found when try to get it from request",
		})
	}

	res, code := c.AppService.CurrentUser(email)

	ctx.IndentedJSON(code, res)
}
