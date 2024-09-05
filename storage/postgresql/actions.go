package postgresql

import (
	"errors"
	"social_media_backend/storage/models"
)

func (s *Storage) IsExistByEmail(email string) bool {
	var userFromDB models.User
	err := s.DB.First(&userFromDB, models.User{Email: email}).Error

	return err == nil && userFromDB.Email == email
}

func (s *Storage) GetUserByEmail(email string) (models.User, error) {
	var userFromDB models.User
	err := s.DB.First(&userFromDB, models.User{Email: email}).Error
	if err != nil {
		return models.User{}, errors.New("User not found")
	}
	return userFromDB, nil
}
