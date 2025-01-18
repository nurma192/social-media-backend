package routes

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"social-media-back/config"
	"social-media-back/controllers"
	"social-media-back/internal/auth"
	"social-media-back/internal/awsStorage"
	middleware "social-media-back/internal/middlware"
	"social-media-back/internal/redisStorage"
	"social-media-back/services"
)

func SetupRoutes(config *config.Config, db *sql.DB, redisClient *redis.Client) *gin.Engine {
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())

	redisService := redisStorage.NewRedisService(redisClient)
	jwtService := auth.NewJWTService(config)
	awsService := awsStorage.InitAWS(&awsStorage.MyConfig{
		AWSRegion:      config.AWSRegion,
		AWSAccessKeyID: config.AWSAccessKeyID,
		AWSSecretKey:   config.AWSSecretKey,
		AWSS3Bucket:    config.AWSS3Bucket,
	})
	appService := services.NewAppService(db, jwtService, redisService, awsService)
	appController := controllers.NewController(appService)

	router.Static("uploads", "./uploads")

	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", appController.Login)
		authGroup.POST("/register", appController.Register)
		authGroup.POST("/send-verify-code", appController.SendVerifyCode)
		authGroup.POST("/verify-account", appController.VerifyAccount)
		authGroup.POST("/refresh", appController.RefreshToken)
	}

	userGroup := router.Group("/user").Use(appService.RequireAuth)
	{
		userGroup.GET("/current", appController.Current)
	}

	postGroup := router.Group("/posts") //.Use(appService.RequireAuth)
	{
		postGroup.POST("", appController.CreatePost)
		postGroup.GET("/:id", appController.GetPost)
		postGroup.GET("", appController.GetAllPosts)
		postGroup.DELETE("/:id", appController.DeletePost)
		//postGroup.PUT("/:id", appController.DeletePost2)
	}

	return router
}
