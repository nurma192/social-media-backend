package main

import (
	"social-media-back/config"
	"social-media-back/internal/logger"
	"social-media-back/internal/storage"
	"social-media-back/routes"
)

func main() {
	log := logger.NewLogger()

	cfg := config.LoadConfig()

	log.Info("Config: ", "connection", cfg)

	db, err := storage.ConnectDB(cfg)
	if err != nil {
		log.Error("Error when try to connect to database", "Error", err)
		return
	}

	defer db.Close()

	router := routes.SetupRoutes(db)

	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Error("Error when starting the server", "Error", err)
		return
	}

	log.Info("Application started successfully!")
}
