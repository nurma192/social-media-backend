package postController

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
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

		currentUser := actions.GetUserFromReq(c)
		if currentUser.ID == 0 {
			return
		}

		var post models.Post
		if err := storage.DB.Limit(1).Find(&post, &models.Post{ID: id, UserID: currentUser.ID}).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error from server": err.Error()})
			return
		}
		if post.ID == 0 {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Post not found"})
			return
		}

		if err := storage.DB.Delete(&post).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error when try to delete the post: " + err.Error()})
			log.Error(err.Error())
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{
			"message": "Post successfully deleted",
			"post":    post,
		})

	}
}
