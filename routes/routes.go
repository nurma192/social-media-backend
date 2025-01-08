package routes

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"social-media-back/controllers"
	"social-media-back/services"
)

func SetupRoutes(db *sql.DB, redisClient *redis.Client) *gin.Engine {
	router := gin.Default()

	appService := services.NewAppService(db, redisClient)
	appController := controllers.NewController(appService)

	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", appController.Login)
		authGroup.POST("/register", appController.Register)
		authGroup.POST("/send-verify-code", appController.SendVerifyCode)
		authGroup.POST("/verify-account", appController.VerifyAccount)
	}

	return router
}
