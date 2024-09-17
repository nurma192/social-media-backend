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
	err := s.DB.Limit(1).Find(&userFromDB, userDetails).Error
	if err != nil {
		return models.User{}, errors.New("user not found")
	}
	return userFromDB, nil
}

func (s *Storage) IsThisUserFollowedTo(userID, checkUserID uint) bool {
	var user models.User
	err := s.DB.Preload("Followings").First(&user, models.User{ID: userID}).Error
	if err != nil {
		return false
	}

	for _, user := range user.Followings {
		if user.ID == checkUserID {
			return true
		}
	}
	return false
}

func (s *Storage) GetPostByID(id uint) (*models.Post, error) {
	var post models.Post
	err := s.DB.First(&post, models.Post{ID: id}).Error
	if err != nil {
		return nil, errors.New("post not found")
	}
	return &post, nil
}
