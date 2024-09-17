package postController

import (
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
	"gorm.io/gorm"
	"net/http"
	"social_media_backend/lib/actions"
	"social_media_backend/storage/models"
	"social_media_backend/storage/postgresql"
)

func DeletePost(log *slog.Logger, storage *postgresql.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := actions.ParseToUint(c.Param("id"))
		if !ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid post id"})
			return
		}

		object, _ := c.Get("user")
		currentUser, ok := object.(models.User)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "User not found in request context"})
			log.Error("User not found in request context")
			return
		}

		user := models.User{
			ID: currentUser.ID,
		}
		err := storage.DB.Preload("Posts").First(&user, user).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			} else {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				log.Error("Failed in func delete post, some error", "error", err)
				return
			}
		}

		var deletePost models.Post
		for _, post := range user.Posts {
			if post.ID == id {
				deletePost = post
				break
			}
		}
		if deletePost.ID == 0 {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}

		err = storage.DB.Delete(&deletePost).Error
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Error(err.Error())
			return
		}
		c.IndentedJSON(http.StatusOK, gin.H{
			"message": "Post successfully deleted",
		})

	}
}
