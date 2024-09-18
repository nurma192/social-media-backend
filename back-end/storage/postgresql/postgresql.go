package postgresql

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"social_media_backend/storage/models"
)

type Storage struct {
	DB *gorm.DB
}

// docker run -it --name social-media-postgres -e POSTGRES_USER=user -e POSTGRES_PASSWORD=pass -e POSTGRES_DB=sm -p 5432:5432 postgres

func NewStorage() (*Storage, error) {
	const op = "storage.postgres.New"
	dsn := "host=localhost user=postgres password=uk888888 dbname=sm port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.Like{},
		&models.Post{},
		&models.Comment{},
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	storage := &Storage{
		DB: db,
	}

	return storage, nil
}
