package postgresql

import "social_media_backend/storage/models"

func (s *Storage) IsExistByEmail(email string) bool {
	var userFromDB models.User
	err := s.DB.First(&userFromDB, models.User{Email: email}).Error

	return err == nil && userFromDB.Email == email

}
