package likeController

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
	"net/http"
	"social_media_backend/lib/actions"
	"social_media_backend/storage/models"
	"social_media_backend/storage/postgresql"
)

func UnLikePost(log *slog.Logger, storage *postgresql.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("DELETE")
		postID, ok := actions.ParseToUint(c.Param("id"))
		if !ok {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid post id"})
			return
		}

		currentUser := actions.GetUserFromReq(c)
		if currentUser.ID == 0 {
			return
		}

		var like models.Like
		if err := storage.DB.Limit(1).Find(&like, models.Like{PostID: postID, UserID: currentUser.ID}).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Error("Error when try to find is like exist: ", err)
			return
		}

		if like.ID == 0 {
			c.IndentedJSON(http.StatusConflict, gin.H{"error": "The post wasn't liked anyway"})
			return
		}

		if err := storage.DB.Delete(&like).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Error("Error when try to delete like: ", err)
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{"message": "The post has been unLiked"})

	}
}
