package storage

import (
	"database/sql"
	"errors"
	"fmt"
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
		&user.Id, &user.Email, &user.Username, &user.Firstname, &user.Lastname, &user.Password, &user.AvatarURL,
		&user.DateOfBirth, &user.Bio, &user.Verified, &user.Location, &user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no user found") // No user found, but no error
		}
		return nil, err
	}
	return user, nil
}
func (s *DBService) GetUserOnlyMainInfoById(id string) (*models.UserMainInfo, error) {
	user := &models.UserMainInfo{}
	err := s.DB.QueryRow(
		"SELECT id, username, firstname, lastname, avatar_url FROM users WHERE id = $1",
		id,
	).Scan(
		&user.Id, &user.Username, &user.Firstname, &user.Lastname, &user.AvatarURL,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no user found") // No user found, but no error
		}
		return nil, err
	}
	return user, nil
}

func (s *DBService) GetPostsLikesCount(id string) (int, error) {
	var likes int
	query := "SELECT COUNT(*) FROM likes WHERE id = $1"
	err := s.DB.QueryRow(query, id).Scan(&likes)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return likes, nil
}
func (s *DBService) GetPostsCommentsCount(id string) (int, error) {
	var commentsCount int
	query := "SELECT COUNT(*) FROM comments WHERE id = $1"
	err := s.DB.QueryRow(query, id).Scan(&commentsCount)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return commentsCount, nil
}

func (s *DBService) IsUserExistByEmail(email string) (bool, error) {
	var userId int
	err := s.DB.QueryRow(
		"SELECT id FROM users WHERE email = $1",
		email,
	).Scan(&userId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("no user found")
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

func (s *DBService) GetPostQuery(postId string) (*models.Post, error) {
	query := `SELECT id, user_id, content, created_at FROM posts WHERE id = $1`
	var post models.Post
	err := s.DB.QueryRow(query, postId).Scan(&post.Id, &post.UserId, &post.ContentText, &post.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("post not found")
		}
		return nil, err
	}

	postImages, err := s.GetPostImages(postId)
	if err != nil {
		return nil, err
	}
	post.Images = postImages

	return &post, nil
}
func (s *DBService) GetPostWithAllInfo(postId string) (*models.PostWithAllInfo, error) {
	query := `SELECT id, user_id, content, created_at FROM posts WHERE id = $1`
	var post models.Post
	err := s.DB.QueryRow(query, postId).Scan(&post.Id, &post.UserId, &post.ContentText, &post.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("post not found")
		}
		return nil, err
	}

	user, err := s.GetUserOnlyMainInfoById(post.UserId)
	if err != nil {
		return nil, err
	}

	postImages, err := s.GetPostImages(postId)
	if err != nil {
		return nil, err
	}
	post.Images = postImages

	likesCount, err := s.GetPostsLikesCount(postId)
	if err != nil {
		return nil, err
	}
	commentsCount, err := s.GetPostsCommentsCount(postId)
	if err != nil {
		return nil, err
	}

	postWithUser := &models.PostWithAllInfo{
		Id:            post.Id,
		User:          user,
		ContentText:   post.ContentText,
		LikesCount:    likesCount,
		CommentsCount: commentsCount,
		CreatedAt:     post.CreatedAt,
	}

	return postWithUser, nil
}

func (s *DBService) GetPostsUserIdByPostId(postId string) (string, error) {
	getPostQuery := `SELECT user_id FROM posts WHERE id = $1`
	var userId string
	err := s.DB.QueryRow(getPostQuery, postId).Scan(&userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("post not found")
		}
		return "", err
	}

	return userId, nil
}

func (s *DBService) GetPostImages(postId string) ([]models.Image, error) {
	getPostImagesQuery := `SELECT id,image_url FROM postImages WHERE post_id = $1`
	rows, err := s.DB.Query(getPostImagesQuery, postId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []models.Image{}, nil
		}
		return nil, err
	}
	defer rows.Close()

	var images []models.Image
	for rows.Next() {
		var image models.Image
		if err := rows.Scan(&image.Id, &image.Url); err != nil {
			return nil, err
		}
		images = append(images, image)
	}

	return images, nil
}
