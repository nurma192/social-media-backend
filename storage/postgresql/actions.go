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

func (s *Storage) GetUserBy(userDetails models.User) (models.User, error) {
	var userFromDB models.User
	err := s.DB.First(&userFromDB, userDetails).Error
	if err != nil {
		return models.User{}, errors.New("user not found")
	}
	return userFromDB, nil
}
