package DatabaseService

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"social-media-back/models"
)

func (s *DBService) CreatePostComment(content string, postId, userId int) error {
	query := "INSERT INTO comments (content, user_id, post_id) VALUES ($1, $2, $3)"
	_, err := s.DB.Exec(query, content, userId, postId)
	if err != nil {
		return err
	}
	return nil
}
func (s *DBService) DeletePostComment(commentId, userId int) error {
	query := "DELETE FROM comments WHERE id = $1 AND user_id = $2"
	_, err := s.DB.Exec(query, commentId, userId)
	if err != nil {
		return err
	}
	return nil
}
func (s *DBService) UpdatePostComment(commentId int, content string, userId int) error {
	query := "UPDATE comments SET content = $1 WHERE id = $2 AND user_id = $3"
	result, err := s.DB.Exec(query, content, commentId, userId)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("comment not found or user does not have permission")
	}
	return nil
}

func (s *DBService) GetPostsCommentsCount(postId int) (int, error) {
	var commentsCount int
	query := "SELECT COUNT(*) FROM comments WHERE post_id = $1"
	err := s.DB.QueryRow(query, postId).Scan(&commentsCount)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, fmt.Errorf("PostActions.GetPostsCommentsCount: %w", err)
	}
	return commentsCount, nil
}

func (s *DBService) GetPostComments(postId, limit, page int) ([]models.CommentWithUser, int, error) {
	offset := (page - 1) * limit
	var comments []models.CommentWithUser

	query := `
        SELECT 
            c.id, c.content, c.user_id, c.post_id, c.created_at, 
            u.id AS user_id, u.username AS user_username, u.firstname AS user_firstname, u.lastname AS user_lastname, u.avatar_url AS user_avatar_url
        FROM comments c 
        JOIN users u 
        ON c.user_id = u.id 
        WHERE c.post_id = $1 
        LIMIT $2 OFFSET $3
    `
	rows, err := s.DB.Query(query, postId, limit, offset)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return make([]models.CommentWithUser, 0), 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment models.CommentWithUser
		comment.User = &models.UserMainInfo{}

		err := rows.Scan(&comment.Id, &comment.Content, &comment.User.Id, &comment.PostId, &comment.CreatedAt,
			&comment.User.Id, &comment.User.Username, &comment.User.Firstname, &comment.User.Lastname, &comment.User.AvatarURL)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return make([]models.CommentWithUser, 0), 0, err
		}

		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error processing rows: %v", err)
		return make([]models.CommentWithUser, 0), 0, err
	}

	var totalComments int
	commentCountQuery := "SELECT COUNT(*) FROM comments WHERE post_id = $1"
	err = s.DB.QueryRow(commentCountQuery, postId).Scan(&totalComments)

	totalPage := (totalComments + limit - 1) / limit

	if err != nil {
		return make([]models.CommentWithUser, 0), 0, err
	}

	return comments, totalPage, nil
}
