package main

import (
	"fmt"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/exp/slog"
	"os"
	"social_media_backend/controllers"
	"social_media_backend/http-server/getUserByID"
	"social_media_backend/http-server/login"
	"social_media_backend/http-server/register"
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

	authGroup := router.Group("/api")

	{
		authGroup.POST("/register", register.New(log, storage))
		authGroup.POST("/login", login.New(log, storage))
		authGroup.GET("/current", middleware.RequireAuth(log, storage), controllers.Current)
		authGroup.GET("/users/:id", middleware.RequireAuth(log, storage), getUserByID.New(log, storage))
		authGroup.PUT("/users/:id", middleware.RequireAuth(log, storage), controllers.UpdateUser)

	}

	router.Run(":8092")
}

func setupLogger() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{Level: slog.LevelDebug},
	}

	handler := opts.NewPrettyHandler(os.Stdout)
	return slog.New(handler)
}
