package DatabaseService

import (
	"database/sql"
	"errors"
	"social-media-back/models"
)

func (s *DBService) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := s.DB.QueryRow(
		"SELECT id, email, username, firstname, lastname, password, avatar_url, date_of_birth, bio, verified, location, created_at FROM users WHERE email = $1",
		email,
	).Scan(
		&user.Id, &user.Email, &user.Username, &user.Firstname, &user.Lastname, &user.Password, &user.AvatarURL,
		&user.DateOfBirth, &user.Bio, &user.Verified, &user.Location, &user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}
func (s *DBService) GetUserById(id int) (*models.User, error) {
	user := &models.User{}
	err := s.DB.QueryRow(
		"SELECT id, email, username, firstname, lastname, password, avatar_url, date_of_birth, bio, verified, location, created_at FROM users WHERE id = $1",
		id,
	).Scan(
		&user.Id, &user.Email, &user.Username, &user.Firstname, &user.Lastname, &user.Password, &user.AvatarURL,
		&user.DateOfBirth, &user.Bio, &user.Verified, &user.Location, &user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}
func (s *DBService) GetUserOnlyMainInfoById(id int) (*models.UserMainInfo, error) {
	user := &models.UserMainInfo{}
	err := s.DB.QueryRow(
		"SELECT id, username, firstname, lastname, avatar_url FROM users WHERE id = $1",
		id,
	).Scan(
		&user.Id, &user.Username, &user.Firstname, &user.Lastname, &user.AvatarURL,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (s *DBService) IsUserExistByEmail(email string) (bool, error) {
	var userId int
	err := s.DB.QueryRow(
		"SELECT id FROM users WHERE email = $1",
		email,
	).Scan(&userId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *DBService) IsUserExistByUsername(username string) (bool, error) {
	var userId int
	err := s.DB.QueryRow(
		"SELECT id FROM users WHERE username = $1",
		username,
	).Scan(&userId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil // No user found
		}
		return false, err
	}
	return true, nil
}
