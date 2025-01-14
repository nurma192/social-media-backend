package services

import (
	"database/sql"
	"errors"
	"social-media-back/internal/auth"
	"social-media-back/internal/redisStorage"
	"social-media-back/models"
)

type AppService struct {
	DB           *sql.DB
	JWTService   *auth.JWTService
	RedisService *redisStorage.RedisService
}

func NewAppService(db *sql.DB, jwtService *auth.JWTService, redisService *redisStorage.RedisService) *AppService {
	return &AppService{
		DB:           db,
		JWTService:   jwtService,
		RedisService: redisService,
	}
}

func (s *AppService) getUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := s.DB.QueryRow(
		"SELECT id, email, username, firstname, lastname, password, avatar_url, date_of_birth, bio, verified, location, created_at FROM users WHERE email = $1",
		email,
	).Scan(
		&user.ID, &user.Email, &user.Username, &user.Firstname, &user.Lastname, &user.Password, &user.AvatarURL,
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
		return false, err
	}
	return true, nil
}

func (s *AppService) isUserExistByUsername(username string) (bool, error) {
	var userID int
	err := s.DB.QueryRow(
		"SELECT id FROM users WHERE username = $1",
		username,
	).Scan(&userID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil // No user found
		}
		return false, err
	}
	return true, nil
}
