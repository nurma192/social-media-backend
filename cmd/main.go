package main

import (
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
	"net/http"
	"os"
	"social_media_backend/lib/logger/slogpretty"
)

func main() {
	// todo init log
	log := setupLogger()
	_ = log

	// todo init DB

	// todo init Router

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
