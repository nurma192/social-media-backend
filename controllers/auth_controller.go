package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"social-media-back/services"
)

func Login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": services.Login("User", "password")})
}

func Register(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": services.Register("User", "email@example.com")})
}
