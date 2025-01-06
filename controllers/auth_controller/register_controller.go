package auth_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"social-media-back/services/auth_service"
)

func Register(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": auth_service.Register("User", "email@example.com")})
}
