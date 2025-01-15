package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *AppController) Current(ctx *gin.Context) {
	email := ctx.MustGet("email").(string)

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"message": "Hello currentUser " + email,
	})
}
