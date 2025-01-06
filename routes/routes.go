package routes

import (
	"github.com/gin-gonic/gin"
	"social-media-back/controllers"
)

func SetupRoutes() *gin.Engine {
	router := gin.Default()

	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", controllers.Login)
		authGroup.POST("/register", controllers.Register)
	}

	return router
}
