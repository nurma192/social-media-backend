package services

import (
	"fmt"
	"net/http"
	"social-media-back/models"
	"social-media-back/models/request"
	"social-media-back/models/response"
	"time"
)

func (s *AppService) CreatePost(post request.CreatePostRequest) (*response.CreatePostResponse, int, *response.DefaultResponse) {
	var uploadedImagesURLs []string

	for _, fileHeader := range post.Images {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, http.StatusInternalServerError, &response.DefaultResponse{
				Message: "Failed to open image",
				Detail:  err.Error(),
			}
		}
		defer file.Close()

		// todo create new unique filename
		fileURL, err := s.AWSService.UploadFile(file, fileHeader.Filename, fileHeader.Header.Get("Content-Type"))
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

	var postID int
	err := s.DB.QueryRow(createPostQuery, post.UserID, post.ContentText).Scan(&postID)
	if err != nil {
		return nil, http.StatusInternalServerError, &response.DefaultResponse{
			Message: "Failed to create post on DB",
			Detail:  err.Error(),
		}
	}
	fmt.Println("postID:", postID)

	createPostImagesQuery := `INSERT INTO postImages (post_id, image_url) VALUES ($1, $2)`
	for _, imageURL := range uploadedImagesURLs {
		_, err = s.DB.Exec(createPostImagesQuery, postID, imageURL)
	}

	return &response.CreatePostResponse{
		Message: "Post created successfully",
		Success: true,
		Post: &models.Post{
			ID:          postID,
			ContentText: post.ContentText,
			Images:      uploadedImagesURLs,
			CreatedAt:   time.Now(),
		},
	}, http.StatusCreated, nil
}
