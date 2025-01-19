package storage

import (
	"database/sql"
	"errors"
	"social-media-back/models"
)

type DBService struct {
	DB *sql.DB
}

func NewDBService(db *sql.DB) *DBService {
	return &DBService{DB: db}
}

func (s *DBService) GetUserByEmail(email string) (*models.User, error) {
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

func (s *DBService) IsUserExistByEmail(email string) (bool, error) {
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

func (s *DBService) IsUserExistByUsername(username string) (bool, error) {
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
