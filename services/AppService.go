package services

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-redis/redis/v8"
	"social-media-back/models"
)

type AppService struct {
	DB          *sql.DB
	RedisClient *redis.Client
	RedisCtx    context.Context
}

func NewAppService(db *sql.DB, redisClient *redis.Client) *AppService {
	return &AppService{
		DB:          db,
		RedisClient: redisClient,
		RedisCtx:    context.Background(),
	}
}
func (s *AppService) isUserExistByEmail(email string) (bool, error) {
	var userID int
	err := s.DB.QueryRow(
		"SELECT id FROM users WHERE email = $1",
		email,
	).Scan(&userID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil // No user found
		}
		return false, err // Some other error occurred
	}
	return true, nil // User exists
}

func (s *AppService) getUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := s.DB.QueryRow(
		"SELECT id, email, firstname, lastname, avatar_url, date_of_birth, bio, verified, location, created_at FROM users WHERE email = $1",
		email,
	).Scan(
		&user.ID, &user.Email, &user.Firstname, &user.Lastname, &user.AvatarURL,
		&user.DateOfBirth, &user.Bio, &user.Verified, &user.Location, &user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // No user found, but no error
		}
		return nil, err
	}
	return user, nil
}
