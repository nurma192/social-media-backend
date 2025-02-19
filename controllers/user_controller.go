package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"social-media-back/models/response"
	"strconv"
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

func (c *AppController) GetUserById(ctx *gin.Context) {
	id := ctx.Param("id")
	fmt.Println(id)

	userId, err := strconv.Atoi(id)
	if err != nil {
		res := response.Response{
			Error: err.Error(),
			Code:  http.StatusBadRequest,
		}
		ctx.IndentedJSON(res.Code, res)
		return
	}

	res := c.AppService.GetUserById(userId)
	ctx.IndentedJSON(res.Code, res)
}
