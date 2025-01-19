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

	createPostQuery := `INSERT INTO posts (user_id, content) VALUES ($1, $2) RETURNING id`
	var postID string
	err := s.DBService.DB.QueryRow(createPostQuery, userId, post.ContentText).Scan(&postID)
	if err != nil {
		return nil, http.StatusInternalServerError, &response.DefaultResponse{
			Message: "Failed to create post on DB",
			Detail:  err.Error(),
		}
	}

	createPostImagesQuery := `INSERT INTO postImages (post_id, image_url) VALUES ($1, $2) RETURNING id`
	var imageID string
	var images []models.Image
	for _, imageURL := range uploadedImagesURLs {
		err = s.DBService.DB.QueryRow(createPostImagesQuery, postID, imageURL).Scan(&imageID)
		if err != nil {
			return nil, http.StatusInternalServerError, &response.DefaultResponse{
				Message: "Failed to create post images on DB",
				Detail:  err.Error(),
			}
		}
		images = append(images, models.Image{
			Id:  imageID,
			Url: imageURL,
		})
	}

	return &response.CreatePostResponse{
		Message: "Post created successfully",
		Success: true,
		Post: &models.Post{
			Id:          postID,
			ContentText: post.ContentText,
			Images:      images,
			UserId:      userId,
			CreatedAt:   time.Now(),
		},
	}, http.StatusCreated, nil
}

func (s *AppService) GetPostById(postID string) (*models.Post, int, *response.DefaultResponse) {
	getPostQuery := `SELECT id, content, created_at FROM posts WHERE id = $1`
	var post models.Post
	err := s.DBService.DB.QueryRow(getPostQuery, postID).Scan(&post.Id, &post.ContentText, &post.CreatedAt)
	if err != nil {
		return nil, http.StatusInternalServerError, &response.DefaultResponse{
			Message: "Failed to get post from DB",
			Detail:  err.Error(),
		}
	}

	postImages, err := s.DBService.GetPostImages(postID)
	if err != nil {
		return nil, http.StatusInternalServerError, &response.DefaultResponse{
			Message: "Failed to get post images from DB",
			Detail:  err.Error(),
		}
	}
	post.Images = postImages

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
		postImages, err := s.DBService.GetPostImages(postID)
		if err != nil {
			return nil, http.StatusInternalServerError, &response.DefaultResponse{
				Message: "Failed to get post images from DB",
				Detail:  err.Error(),
			}
		}
		posts = append(posts, &models.Post{
			Id:          postID,
			ContentText: content,
			UserId:      userID,
			Images:      postImages,
			CreatedAt:   createdAt,
		})
	}

	res := &response.GetPostsResponse{
		Posts:   posts,
		Success: true,
		Page:    page,
		Limit:   limit,
	}
	return res, http.StatusOK, nil
}

func (s *AppService) DeletePost(postID, userId string) (*response.DefaultResponse, int, *response.DefaultResponse) {
	postsUserId, err := s.DBService.GetPostsUserIdByPostId(postID)
	if err != nil {
		return nil, http.StatusInternalServerError, &response.DefaultResponse{
			Message: err.Error(),
			Detail:  "Failed to get post user id from DB",
		}
	}

	if postsUserId != userId {
		return nil, http.StatusForbidden, &response.DefaultResponse{
			Message: "You are not allowed to delete this post",
		}
	}

	deletePostQuery := `DELETE FROM posts WHERE id = $1`
	_, err = s.DBService.DB.Exec(deletePostQuery, postID)
	if err != nil {
		return nil, http.StatusInternalServerError, &response.DefaultResponse{
			Message: "Failed to delete post from DB",
			Detail:  err.Error(),
		}
	}

	return &response.DefaultResponse{
		Message: "Post deleted successfully",
		Success: true,
	}, http.StatusOK, nil
}

func (s *AppService) UpdatePost(postId, userId string, req *request.UpdatePostRequest) (*response.UpdatePostResponse, int, *response.DefaultResponse) {
	post, err := s.DBService.GetPostQuery(postId)
	if err != nil {
		return nil, http.StatusInternalServerError, &response.DefaultResponse{
			Message: err.Error(),
			Detail:  "Failed to get post from DB",
		}
	}

	if len(post.Images)-len(req.DeletedImagesId)+len(req.NewImages) > 5 {
		return nil, http.StatusBadRequest, &response.DefaultResponse{
			Message: "You can't add more than 5 images to post",
		}
	}

	if post.UserId != userId {
		return nil, http.StatusForbidden, &response.DefaultResponse{
			Message: "You are not allowed to update this post",
		}
	}

	for _, deleteImageId := range req.DeletedImagesId {
		var deleteImageURL *string
		for _, image := range post.Images {
			if image.Id == deleteImageId {
				deleteImageURL = &image.Url
				break
			}
		}
		if deleteImageURL == nil {
			return nil, http.StatusForbidden, &response.DefaultResponse{
				Message: "You are not allowed to delete this image",
			}
		}

		err := s.AWSService.DeleteFile(*deleteImageURL)
		if err != nil {
			return nil, http.StatusInternalServerError, &response.DefaultResponse{
				Message: "Failed to delete image from Amazon S3",
				Detail:  err.Error(),
			}
		}
	}

	// delete images from DB
	deleteImageQuery := `DELETE FROM postImages WHERE id = $1`
	for _, imageId := range req.DeletedImagesId {
		_, err = s.DBService.DB.Exec(deleteImageQuery, imageId)
		if err != nil {
			return nil, http.StatusInternalServerError, &response.DefaultResponse{
				Message: "Failed to delete image from DB",
				Detail:  err.Error(),
			}
		}
	}

	// upload images to S3 Storage
	var uploadedImagesURLs []string
	for index, fileHeader := range req.NewImages {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, http.StatusInternalServerError, &response.DefaultResponse{
				Message: "Failed to open image",
				Detail:  err.Error(),
			}
		}
		defer file.Close()

		fileURL, err := s.AWSService.UploadFile(file, fileHeader.Filename+time.Now().Format("_2006-01-02_15:04:05")+fmt.Sprintf("%d", index)+userId, fileHeader.Header.Get("Content-Type"))
		uploadedImagesURLs = append(uploadedImagesURLs, fileURL)
	}

	// add new postImages
	addNewPostImagesQuery := `INSERT INTO postImages (post_id, image_url) VALUES ($1, $2)`
	for _, imageURL := range uploadedImagesURLs {
		_, err := s.DBService.DB.Exec(addNewPostImagesQuery, postId, imageURL)
		if err != nil {
			return nil, http.StatusInternalServerError, &response.DefaultResponse{
				Message: "Failed to add new post images to DB",
				Detail:  err.Error(),
			}
		}
	}

	// update post content
	updatePostQuery := `UPDATE posts SET content = $1 WHERE id = $2`
	_, err = s.DBService.DB.Exec(updatePostQuery, req.ContentText, postId)
	if err != nil {
		return nil, http.StatusInternalServerError, &response.DefaultResponse{
			Message: "Failed to update post on DB",
			Detail:  err.Error(),
		}
	}

	post, err = s.DBService.GetPostQuery(postId)

	return &response.UpdatePostResponse{
		Success: true,
		Message: "Post updated successfully",
		Post:    post,
	}, http.StatusOK, nil
}
