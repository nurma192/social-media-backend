package DatabaseService

import (
	"database/sql"
	"errors"
	"fmt"
)

// pgx
func (s *DBService) AddLikePost(postId, userId int) error {
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

func (s *DBService) DeleteLikePost(postId, userId int) error {
	deletePostQuery := `DELETE FROM likes WHERE post_id = $1 AND user_id = $2`
	result, err := s.DB.Exec(deletePostQuery, postId, userId)
	if err != nil {
		return fmt.Errorf("LikeActions.DeleteLikePost: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("post not liked by user")
	}

	return nil
}

func (s *DBService) IsUserLikedPost(postId, userId int) (bool, error) {
	query := `SELECT 1 FROM likes WHERE post_id = $1 AND user_id = $2`
	var exists int
	err := s.DB.QueryRow(query, postId, userId).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("LikeActions.IsUserLikedPost: %w", err)
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
		return 0, fmt.Errorf("LikeActions.GetPostsLikesCount: %w", err)
	}
	return likes, nil
}
