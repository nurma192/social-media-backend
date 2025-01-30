package services

import (
	"fmt"
	"net/http"
	"social-media-back/models"
	"social-media-back/models/request"
	"social-media-back/models/response"
	"strconv"
	"time"
)

func (s *AppService) CreatePost(post request.CreatePostRequest, userId int) (*response.Response, int) {
	var uploadedImagesURLs []string

	for index, fileHeader := range post.Images {
		file, err := fileHeader.Open()
		if err != nil {
			return &response.Response{
				Error: err.Error(),
			}, http.StatusInternalServerError
		}
		defer file.Close()

		// todo create new unique filename
		fileURL, err := s.AWSService.UploadFile(
			file,
			fileHeader.Filename+time.Now().Format("_2006-01-02_15:04:05")+fmt.Sprintf("%d", index)+strconv.Itoa(userId),
			fileHeader.Header.Get("Content-Type"),
		)

		if err != nil {
			return &response.Response{
				Error: err.Error(),
			}, http.StatusInternalServerError
		}
		uploadedImagesURLs = append(uploadedImagesURLs, fileURL)
	}

	createPostQuery := `INSERT INTO posts (user_id, content) VALUES ($1, $2) RETURNING id`
	fmt.Println(createPostQuery, userId, post.ContentText)
	var postID int
	err := s.DBService.DB.QueryRow(createPostQuery, userId, post.ContentText).Scan(&postID)
	if err != nil {
		return &response.Response{
			Error: err.Error(),
		}, http.StatusInternalServerError
	}

	createPostImagesQuery := `INSERT INTO postImages (post_id, image_url) VALUES ($1, $2) RETURNING id`
	var imageID int
	var images []models.Image
	for _, imageURL := range uploadedImagesURLs {
		err = s.DBService.DB.QueryRow(createPostImagesQuery, postID, imageURL).Scan(&imageID)
		if err != nil {
			return &response.Response{
				Error: err.Error(),
			}, http.StatusInternalServerError
		}
		images = append(images, models.Image{
			Id:  imageID,
			Url: imageURL,
		})
	}

	return &response.Response{
		Result: &models.Post{
			Id:          postID,
			ContentText: post.ContentText,
			Images:      images,
			UserId:      userId,
			CreatedAt:   time.Now(),
		},
		Message: "Post created",
	}, http.StatusCreated
}

func (s *AppService) GetPostById(postId, userId int) (*response.Response, int) {
	postWithAllInfo, err := s.DBService.GetPostWithAllInfo(postId)
	if err != nil {
		return &response.Response{
			Error: err.Error(),
		}, http.StatusInternalServerError
	}

	postImages, err := s.DBService.GetPostImages(postId)
	if err != nil {
		return &response.Response{
			Error: err.Error(),
		}, http.StatusInternalServerError
	}
	postWithAllInfo.Images = postImages

	isUserLikedPost, err := s.DBService.IsUserLikedPost(postId, userId)
	if err != nil {
		postWithAllInfo.LikedByUser = false
	}
	postWithAllInfo.LikedByUser = isUserLikedPost

	return &response.Response{
		Result: postWithAllInfo,
	}, http.StatusOK
}

func (s *AppService) GetAllPosts(limit, page, userId int) (*response.PaginationResponse, int) {
	posts, totalPages, err := s.DBService.GetAllPostsWithAllInfo(limit, page, userId)

	if err != nil {
		return &response.PaginationResponse{
			Error: err.Error(),
		}, http.StatusInternalServerError
	}

	res := &response.PaginationResponse{
		Result:     posts,
		Page:       page,
		TotalPages: totalPages,
		Limit:      limit,
	}
	return res, http.StatusOK
}

func (s *AppService) DeletePost(postID, userId int) (*response.Response, int) {
	postsUserId, err := s.DBService.GetPostsUserIdByPostId(postID)
	if err != nil {
		return &response.Response{
			Error: err.Error(),
		}, http.StatusInternalServerError
	}

	if postsUserId != userId {
		return &response.Response{
			Error: "You are not allowed to delete this post",
		}, http.StatusForbidden
	}

	deletePostQuery := `DELETE FROM posts WHERE id = $1`
	_, err = s.DBService.DB.Exec(deletePostQuery, postID)
	if err != nil {
		return &response.Response{
			Error: err.Error(),
		}, http.StatusInternalServerError
	}

	return &response.Response{
		Message: "Post deleted successfully",
	}, http.StatusOK
}

func (s *AppService) UpdatePost(postId, userId int, req *request.UpdatePostRequest) (*response.Response, int) {
	post, err := s.DBService.GetPostQuery(postId)
	if err != nil {
		return &response.Response{
			Error: err.Error(),
		}, http.StatusInternalServerError
	}

	if len(post.Images)-len(req.DeletedImagesId)+len(req.NewImages) > 5 {
		return &response.Response{
			Error: "You can't add more than 5 images to post",
		}, http.StatusBadRequest
	}

	if post.UserId != userId {
		return &response.Response{
			Error: "You are not allowed to update this post",
		}, http.StatusForbidden
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
			return &response.Response{
				Error: "You are not allowed to delete this image",
			}, http.StatusForbidden
		}

		err := s.AWSService.DeleteFile(*deleteImageURL)
		if err != nil {
			return &response.Response{
				Error: err.Error(),
			}, http.StatusInternalServerError
		}
	}

	// delete images from DB
	deleteImageQuery := `DELETE FROM postImages WHERE id = $1`
	for _, imageId := range req.DeletedImagesId {
		_, err = s.DBService.DB.Exec(deleteImageQuery, imageId)
		if err != nil {
			return &response.Response{
				Error: err.Error(),
			}, http.StatusInternalServerError
		}
	}

	// upload images to S3 Storage
	var uploadedImagesURLs []string
	for index, fileHeader := range req.NewImages {
		file, err := fileHeader.Open()
		if err != nil {
			return &response.Response{
				Error: err.Error(),
			}, http.StatusInternalServerError
		}
		defer file.Close()

		fileURL, err := s.AWSService.UploadFile(file, fileHeader.Filename+time.Now().Format("_2006-01-02_15:04:05")+fmt.Sprintf("%d", index)+strconv.Itoa(userId), fileHeader.Header.Get("Content-Type"))
		uploadedImagesURLs = append(uploadedImagesURLs, fileURL)
	}

	// add new postImages
	addNewPostImagesQuery := `INSERT INTO postImages (post_id, image_url) VALUES ($1, $2)`
	for _, imageURL := range uploadedImagesURLs {
		_, err := s.DBService.DB.Exec(addNewPostImagesQuery, postId, imageURL)
		if err != nil {
			return &response.Response{
				Error: err.Error(),
			}, http.StatusInternalServerError
		}
	}

	// update post content
	updatePostQuery := `UPDATE posts SET content = $1 WHERE id = $2`
	_, err = s.DBService.DB.Exec(updatePostQuery, req.ContentText, postId)
	if err != nil {
		return &response.Response{
			Error: err.Error(),
		}, http.StatusInternalServerError
	}

	post, err = s.DBService.GetPostQuery(postId)

	return &response.Response{
		Message: "Post updated successfully",
		Result:  post,
	}, http.StatusOK
}
