package main

import (
	"fmt"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"golang.org/x/exp/slog"
	"net/http"
	"os"
	"social_media_backend/lib/logger/slogpretty"
	"social_media_backend/storage/postgresql"
)

func main() {
	// todo init log
	log := setupLogger()
	_ = log

	// todo init DB

	storage, err := postgresql.NewStorage()
	if err != nil {
		log.Error("Error creating postgres storage", err.Error())
		os.Exit(1)
	}
	log.Info("Successfully connected to postgres storage")
	fmt.Println(storage)

	// todo init Router

	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	router.Use(requestid.New())

	authGroup := router.Group("/auth")

	{
		authGroup.GET("/login", func(c *gin.Context) {
			c.IndentedJSON(http.StatusOK, gin.H{
				"message": "login success",
			})
		})
		authGroup.GET("/register", func(c *gin.Context) {
			c.IndentedJSON(http.StatusOK, gin.H{
				"message": "register success",
			})
		})
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
