package postController

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
	"net/http"
	"social_media_backend/lib/actions"
	"social_media_backend/storage/models"
	"social_media_backend/storage/postgresql"
)

func GetPostByID(log *slog.Logger, storage *postgresql.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		currentUser := actions.GetUserFromReq(c)
		if currentUser.ID == 0 {
			return
		}

		var post models.Post
		if err := storage.DB.Preload("Likes").Preload("Comments").Limit(1).Find(&post, id).Error; err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Error("Failed to get post by ID", "error", err)
			return
		}

		if post.ID == 0 {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "post not found"})
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
