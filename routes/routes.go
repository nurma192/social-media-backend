package routes

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"social-media-back/config"
	"social-media-back/controllers"
	"social-media-back/internal/awsStorage"
	"social-media-back/internal/mail"
	middleware "social-media-back/internal/middlware"
	"social-media-back/internal/redisStorage"
	"social-media-back/internal/storage/DatabaseService"
	"social-media-back/internal/token"
	"social-media-back/services"
)

func SetupRoutes(config *config.Config, db *sql.DB, redisClient *redis.Client) *gin.Engine {
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())

	databaseService := DatabaseService.NewDBService(db)
	redisService := redisStorage.NewRedisService(redisClient)
	jwtService := token.NewJWTService(config)
	awsService := awsStorage.InitAWS(&awsStorage.MyConfig{
		AWSRegion:      config.AWSRegion,
		AWSAccessKeyID: config.AWSAccessKeyID,
		AWSSecretKey:   config.AWSSecretKey,
		AWSS3Bucket:    config.AWSS3Bucket,
	})
	emailService := mail.NewEmailService()

	appService := services.NewAppService(databaseService, jwtService, redisService, awsService, emailService)
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
		userGroup.GET("/:id", appController.GetUserById)
	}

	postGroup := router.Group("/posts").Use(appService.RequireAuth)
	{
		postGroup.POST("", appController.CreatePost)
		postGroup.GET("/:id", appController.GetPost)
		postGroup.GET("", appController.GetAllPosts)
		postGroup.DELETE("/:id", appController.DeletePost)
		postGroup.PUT("/:id", appController.UpdatePost)
	}

	likeGroup := router.Group("/like").Use(appService.RequireAuth)
	{
		likeGroup.POST("", appController.LikePost)
		likeGroup.DELETE("", appController.UnlikePost)
	}

	commentGroup := router.Group("/postComments").Use(appService.RequireAuth)
	{
		commentGroup.POST("", appController.CreatePostComment)
		commentGroup.DELETE("", appController.DeletePostComment)
		commentGroup.PUT("", appController.UpdatePostComment)
		commentGroup.GET("/:id", appController.GetPostsComments)
	}

	return router
}
