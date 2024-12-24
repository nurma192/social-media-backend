package commentController

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
	"net/http"
	"social_media_backend/lib/actions"
	"social_media_backend/storage/models"
	"social_media_backend/storage/postgresql"
)

func DeleteComment(log *slog.Logger, storage *postgresql.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := actions.ParseToUint(c.Param("id"))
		if !ok {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "id is required"})
			return
		}

		currentUser := actions.GetUserFromReq(c)
		if currentUser.ID == 0 {
			return
		}

		var comment models.Comment
		if err := storage.DB.Limit(1).Find(&comment, &models.Comment{ID: id, UserID: currentUser.ID}).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Error("Error while finding comment by id:", err)
			return
		}

		if comment.ID == 0 {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
			return
		}
		if err := storage.DB.Delete(&comment).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Error("Error while deleting comment by id:", err)
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{
			"comment": comment,
			"message": "Comment deleted",
		})
	}
}
