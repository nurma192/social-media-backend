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

func GetPostByID(log *slog.Logger, storage *postgresql.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		object, _ := c.Get("user")
		currentUser, ok := object.(models.User)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "User not found in request context"})
			log.Error("User not found in request context")
			return
		}

		var post models.Post
		err := storage.DB.Preload("Likes").Preload("Comments").Limit(1).Find(&post, id).Error
		if err != nil || post.ID == 0 {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			} else {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				log.Error("Failed to get post by ID", "error", err)
			}
			return
		}

		isUserLiked := false
		for _, like := range post.Likes {
			if like.UserID == currentUser.ID {
				isUserLiked = true
				break
			}
		}

		c.IndentedJSON(http.StatusOK, gin.H{
			"post":          post,
			"isLikedByUser": isUserLiked,
		})

	}
}
