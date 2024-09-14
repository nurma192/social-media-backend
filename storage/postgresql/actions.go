package postgresql

import (
	"errors"
	"social_media_backend/storage/models"
)

func (s *Storage) IsExistByEmail(email string) bool {
	var userFromDB models.User
	s.DB.Limit(1).Find(&userFromDB, "email = ?", email)

	if userFromDB.Email == email {
		return true
	}

	return false
}

func (s *Storage) GetUserBy(userDetails models.User) (models.User, error) {
	var userFromDB models.User
	err := s.DB.First(&userFromDB, userDetails).Error
	if err != nil {
		return models.User{}, errors.New("user not found")
	}
	return userFromDB, nil
}
