package postgresql

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"social_media_backend/storage/models"
)

type Storage struct {
	db *gorm.DB
}

//docker run -it --name social-media-postgres -e POSTGRES_USER=user -e POSTGRES_PASSWORD=pass -e POSTGRES_DB=sm -p 5432:5432 postgres

func NewStorage() (*Storage, error) {
	const op = "storage.postgres.New"
	db, err := gorm.Open("postgres", "user=user password=pass dbname=sm port=57335 sslmode=disable")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	db.AutoMigrate(&models.Comment{}, &models.Like{}, &models.Post{}, &models.Follow{}, &models.User{})

	storage := &Storage{
		db: db,
	}

	return storage, nil

}
