package services

import (
	"fmt"
	"net/http"
	"social-media-back/models"
	"social-media-back/models/request"
	"social-media-back/models/response"
	"time"
)

func (s *AppService) CreatePost(post request.CreatePostRequest, userId string) (*response.CreatePostResponse, int, *response.DefaultResponse) {
	var uploadedImagesURLs []string

	for index, fileHeader := range post.Images {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, http.StatusInternalServerError, &response.DefaultResponse{
				Message: "Failed to open image",
				Detail:  err.Error(),
			}
		}
		defer file.Close()

		// todo create new unique filename
		fileURL, err := s.AWSService.UploadFile(file, fileHeader.Filename+time.Now().Format("_2006-01-02_15:04:05")+fmt.Sprintf("%d", index)+userId, fileHeader.Header.Get("Content-Type"))
		if err != nil {
			return nil, http.StatusInternalServerError, &response.DefaultResponse{
				Message: "Failed to upload image to Amazon S3",
				Detail:  err.Error(),
			}
		}
		uploadedImagesURLs = append(uploadedImagesURLs, fileURL)
	}
	fmt.Println("uploadedImagesURLs:", uploadedImagesURLs)

	createPostQuery := `INSERT INTO posts (user_id, content) VALUES ($1, $2) RETURNING id`

	var postID string
	err := s.DBService.DB.QueryRow(createPostQuery, userId, post.ContentText).Scan(&postID)
	if err != nil {
		return nil, http.StatusInternalServerError, &response.DefaultResponse{
			Message: "Failed to create post on DB",
			Detail:  err.Error(),
		}
	}
	fmt.Println("postID:", postID)

	createPostImagesQuery := `INSERT INTO postImages (post_id, image_url) VALUES ($1, $2)`
	for _, imageURL := range uploadedImagesURLs {
		_, err = s.DBService.DB.Exec(createPostImagesQuery, postID, imageURL)
	}

	return &response.CreatePostResponse{
		Message: "Post created successfully",
		Success: true,
		Post: &models.Post{
			ID:          postID,
			ContentText: post.ContentText,
			Images:      uploadedImagesURLs,
			UserID:      userId,
			CreatedAt:   time.Now(),
		},
	}, http.StatusCreated, nil
}

func (s *AppService) GetPostById(postID string) (*models.Post, int, *response.DefaultResponse) {
	getPostQuery := `SELECT id, content, created_at FROM posts WHERE id = $1`
	var post models.Post
	err := s.DBService.DB.QueryRow(getPostQuery, postID).Scan(&post.ID, &post.ContentText, &post.CreatedAt)
	if err != nil {
		return nil, http.StatusInternalServerError, &response.DefaultResponse{
			Message: "Failed to get post from DB",
			Detail:  err.Error(),
		}
	}
	getPostImagesQuery := `SELECT image_url FROM postImages WHERE post_id = $1`

	rows, err := s.DBService.DB.Query(getPostImagesQuery, postID)
	if err != nil {
		return nil, http.StatusInternalServerError, &response.DefaultResponse{
			Message: "Failed to get post images from DB",
			Detail:  err.Error(),
		}
	}
	defer rows.Close()

	var images []string
	for rows.Next() {
		var imageURL string
		if err := rows.Scan(&imageURL); err != nil {
			return nil, http.StatusInternalServerError, &response.DefaultResponse{
				Message: "Failed to scan post image URL",
				Detail:  err.Error(),
			}
		}
		images = append(images, imageURL)
	}

	post.Images = images

	return &post, http.StatusOK, nil
}

func (s *AppService) GetAllPosts(limit, page int) (*response.GetPostsResponse, int, *response.DefaultResponse) {
	offset := (page - 1) * limit
	getAllPostsQuery :=
		`SELECT
			id,
			content,
			user_id,
			created_at
		FROM posts
		ORDER BY created_at DESC 
		LIMIT $1 OFFSET $2`

	rows, err := s.DBService.DB.Query(getAllPostsQuery, limit, offset)
	if err != nil {
		return nil, http.StatusInternalServerError, &response.DefaultResponse{
			Message: "Failed to get all posts from DB",
			Detail:  err.Error(),
		}
	}
	defer rows.Close()

	posts := make([]*models.Post, 0)
	for rows.Next() {
		var content, postID, userID string
		var createdAt time.Time

		err := rows.Scan(&postID, &content, &userID, &createdAt)
		if err != nil {
			return nil, http.StatusInternalServerError, &response.DefaultResponse{
				Message: "Failed to scan post data",
				Detail:  err.Error(),
			}
		}
		posts = append(posts, &models.Post{
			ID:          postID,
			ContentText: content,
			UserID:      userID,
			CreatedAt:   createdAt,
		})
	}
	fmt.Println(posts)
	return nil, http.StatusOK, nil
}
