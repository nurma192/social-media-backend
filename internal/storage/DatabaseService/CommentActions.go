package DatabaseService

import "social-media-back/models"

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

func (s *DBService) GetPostComments(postId int) (*[]models.CommentWithUser, error) {
	var comments []models.CommentWithUser
	query := "SELECT id, content, user_id, post_id, created_at FROM comments WHERE post_id = $1"
	rows, err := s.DB.Query(query, postId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment models.CommentWithUser
		err := rows.Scan(&comment.Id, &comment.Content, &comment.User.Id, &comment.PostId, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}

		user, err := s.GetUserOnlyMainInfoById(comment.User.Id)
		if err != nil {
			return nil, err
		}
		comment.User = user
		comments = append(comments, comment)
	}

	return &comments, nil
}
