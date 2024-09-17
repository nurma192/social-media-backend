package postController

import (
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
	"gorm.io/gorm"
	"net/http"
	"social_media_backend/storage/models"
	"social_media_backend/storage/postgresql"
)

type PostWithLikes struct {
	Post          models.Post `json:"post"`
	IsLikedByUser bool        `json:"isLikedByUser"`
}

func GetAllPosts(log *slog.Logger, storage *postgresql.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		object, _ := c.Get("user")

		currentUser, ok := object.(models.User)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "User not found in request context"})
			log.Error("User not found in request context")
			return
		}

		var posts []models.Post
		err := storage.DB.Order("created_at desc").Preload("Likes").Preload("Comments").Find(&posts).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Posts not found"})
			} else {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				log.Error("Failed to get all posts", "error", err)
			}
			return
		}

		if len(posts) == 0 {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Posts not found"})
			return
		}

		var postWithLikes []PostWithLikes
		for _, post := range posts {
			isLiked := false
			for _, like := range post.Likes {
				if like.UserID == currentUser.ID {
					isLiked = true
					continue
				}
			}
			postWithLikes = append(postWithLikes, PostWithLikes{
				Post:          post,
				IsLikedByUser: isLiked,
			})
		}

		c.IndentedJSON(http.StatusOK, postWithLikes)
	}
}
