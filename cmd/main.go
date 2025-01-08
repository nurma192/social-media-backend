package main

import (
	"social-media-back/config"
	"social-media-back/internal/logger"
	"social-media-back/internal/redisStorage"
	"social-media-back/internal/storage"
	"social-media-back/routes"
)

func main() {
	// setup logger
	log := logger.NewLogger()

	// load config
	cfg := config.LoadConfig()
	log.Info("Config: ", "connection", cfg)

	// connect to Postgres
	db, err := storage.ConnectDB(cfg)
	if err != nil {
		log.Error("Error when try to connect to database", "Error", err)
		return
	}
	defer db.Close()

	// connect to Redis
	redisClient := redisStorage.CreateClient()

	// setup router
	router := routes.SetupRoutes(cfg, db, redisClient)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Error("Error when starting the server", "Error", err)
		return
	}

	log.Info("Application started successfully!")
}
