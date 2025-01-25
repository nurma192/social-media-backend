package DatabaseService

import (
	"database/sql"
	"errors"
	"fmt"
)

func (s *DBService) AddLikePost(postId int, userId string) error {
	likedByUser, err := s.IsUserLikedPost(postId, userId)
	if err != nil {
		return err
	}
	if likedByUser {
		return fmt.Errorf("post liked by user")
	}
	getPostQuery := `INSERT INTO likes (post_id, user_id) VALUES ($1, $2)`
	_, err = s.DB.Exec(getPostQuery, postId, userId)
	if err != nil {
		return err
	}
	return nil
}

func (s *DBService) DeleteLikePost(postId int, userId string) error {
	likedByUser, err := s.IsUserLikedPost(postId, userId)
	if err != nil {
		return err
	}
	if !likedByUser {
		return fmt.Errorf("post liked by user")
	}

	deletePostQuery := `DELETE FROM likes WHERE post_id = $1 AND user_id = $2`
	_, err = s.DB.Exec(deletePostQuery, postId, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("post not liked by user")
		}
		return err
	}
	return nil
}

func (s *DBService) IsUserLikedPost(postId int, userId string) (bool, error) {
	query := `SELECT 1 FROM likes WHERE post_id = $1 AND user_id = $2`
	var exists int
	err := s.DB.QueryRow(query, postId, userId).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *DBService) GetPostsLikesCount(postId int) (int, error) {
	var likes int
	query := "SELECT COUNT(*) FROM likes WHERE post_id = $1"
	err := s.DB.QueryRow(query, postId).Scan(&likes)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return likes, nil
}
