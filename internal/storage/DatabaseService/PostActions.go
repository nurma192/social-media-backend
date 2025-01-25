package DatabaseService

import (
	"database/sql"
	"errors"
	"fmt"
	"social-media-back/models"
)

func (s *DBService) GetPostsCommentsCount(id int) (int, error) {
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

func (s *DBService) GetPostQuery(postId int) (*models.Post, error) {
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
func (s *DBService) GetPostWithAllInfo(postId int) (*models.PostWithAllInfo, error) {
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

func (s *DBService) GetPostImages(postId int) ([]models.Image, error) {
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
