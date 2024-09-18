package main

import (
	"fmt"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/exp/slog"
	"os"
	commentController "social_media_backend/http-server/CommentController"
	followController "social_media_backend/http-server/FollowController"
	likeController "social_media_backend/http-server/LikeController"
	postController "social_media_backend/http-server/PostController"
	"social_media_backend/http-server/api/current"
	"social_media_backend/http-server/api/getUserByID"
	"social_media_backend/http-server/api/login"
	"social_media_backend/http-server/api/register"
	"social_media_backend/http-server/api/update"
	"social_media_backend/internal/middleware"
	"social_media_backend/lib/logger/slogpretty"
	"social_media_backend/storage/postgresql"
)

func main() {
	// todo init log
	log := setupLogger()

	//todo init env
	err := godotenv.Load()
	if err != nil {
		log.Error("Error loading .env file")
	}

	// todo init DB
	storage, err := postgresql.NewStorage()
	if err != nil {
		log.Error("Error creating postgres storage", err.Error())
		os.Exit(1)
	}
	log.Info("Successfully connected to postgres storage")
	fmt.Println(storage)

	// todo ------------------------------------------------

	// todo ------------------------------------------------

	// todo init Router
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	router.Use(requestid.New())

	authGroup := router.Group("/auth")

	{
		authGroup.POST("/register", register.New(log, storage))
		authGroup.POST("/login", login.New(log, storage))

	}

	apiGroup := router.Group("/api").Use(middleware.RequireAuth(log, storage))
	{
		apiGroup.GET("/current", current.New(log, storage))
		apiGroup.GET("/users/:id", getUserByID.New(log, storage))
		apiGroup.PUT("/users/:id", update.New(log, storage))
	}

	postGroup := router.Group("/posts").Use(middleware.RequireAuth(log, storage))
	{
		postGroup.POST("/", postController.CreatePost(log, storage))
		postGroup.GET("/", postController.GetAllPosts(log, storage))
		postGroup.GET("/:id", postController.GetPostByID(log, storage))
		postGroup.DELETE("/:id", postController.DeletePost(log, storage))
	}

	commentGroup := router.Group("/comments").Use(middleware.RequireAuth(log, storage))
	{
		commentGroup.POST("/", commentController.CreateComment(log, storage))
		commentGroup.DELETE("/:id", commentController.DeleteComment(log, storage))
	}

	likeGroup := router.Group("/likes").Use(middleware.RequireAuth(log, storage))
	{
		likeGroup.POST("/", likeController.LikePost(log, storage))
		likeGroup.DELETE("/:id", likeController.UnLikePost(log, storage))
	}

	router.POST("/follow", middleware.RequireAuth(log, storage), followController.FollowUser(log, storage))
	router.DELETE("/unfollow/:id", middleware.RequireAuth(log, storage), followController.UnFollowUser(log, storage))

	router.Run(":8092")
}

func setupLogger() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{Level: slog.LevelDebug},
	}

	handler := opts.NewPrettyHandler(os.Stdout)
	return slog.New(handler)
}
