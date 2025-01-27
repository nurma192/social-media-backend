package DatabaseService

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"social-media-back/models"
)

func (s *DBService) CreatePostComment(postId int) error {
	query := "INSERT INTO comments (content, user_id, post_id) VALUES ($1, $2, $3)"
	_, err := s.DB.Exec(query, postId, postId, postId)
	if err != nil {
		return err
	}
	return nil
}
func (s *DBService) DeletePostComment(commentId int) error {
	query := "DELETE FROM comments WHERE id = $1"
	_, err := s.DB.Exec(query, commentId)
	if err != nil {
		return err
	}
	return nil
}
func (s *DBService) UpdatePostComment(commentId int, content string) error {
	query := "UPDATE comments SET content = $1 WHERE id = $2"
	_, err := s.DB.Exec(query, commentId, content)
	if err != nil {
		return err
	}
	return nil
}

func (s *DBService) GetPostsCommentsCount(id int) (int, error) {
	var commentsCount int
	query := "SELECT COUNT(*) FROM comments WHERE id = $1"
	err := s.DB.QueryRow(query, id).Scan(&commentsCount)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, fmt.Errorf("PostActions.GetPostsCommentsCount: %w", err)
	}
	return commentsCount, nil
}

func (s *DBService) GetPostComments(postId, limit, page int) ([]models.CommentWithUser, error) {
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
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment models.CommentWithUser
		comment.User = &models.UserMainInfo{} // Инициализация вложенной структуры

		err := rows.Scan(&comment.Id, &comment.Content, &comment.User.Id, &comment.PostId, &comment.CreatedAt,
			&comment.User.Id, &comment.User.Username, &comment.User.Firstname, &comment.User.Lastname, &comment.User.AvatarURL)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}

		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error processing rows: %v", err)
		return nil, err
	}

	return comments, nil
}
