package routes

import (
	"github.com/gin-gonic/gin"
	"social-media-back/controllers/auth_controller"
)

func SetupRoutes() *gin.Engine {
	router := gin.Default()

	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", auth_controller.Login)
		authGroup.POST("/register", auth_controller.Register)
	}

	return router
}
