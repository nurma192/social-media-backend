package likeController

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
	"net/http"
	"social_media_backend/lib/actions"
	"social_media_backend/storage/models"
	"social_media_backend/storage/postgresql"
)

type LikeReqStruct struct {
	PostID uint
}

func LikePost(log *slog.Logger, storage *postgresql.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var likeRequest LikeReqStruct
		if err := c.ShouldBind(&likeRequest); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if likeRequest.PostID == 0 {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "PostID is required"})
			return
		}

		currentUser := actions.GetUserFromReq(c)
		if currentUser.ID == 0 {
			return
		}

		var like models.Like
		if err := storage.DB.Limit(1).Find(&like, models.Like{PostID: likeRequest.PostID, UserID: currentUser.ID}).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			log.Error("Error when try to find is like exist: ", err)
			return
		}

		if like.ID != 0 {
			c.IndentedJSON(http.StatusConflict, gin.H{"error": "Post is already likely"})
			return
		}

		newLike := models.Like{
			PostID: likeRequest.PostID,
			UserID: currentUser.ID,
		}
		if err := storage.DB.Create(&newLike).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Error("Error when try to create like: ", err)
			return
		}

		c.IndentedJSON(http.StatusCreated, gin.H{
			"like":    newLike,
			"message": "Post liked successfully",
		})

	}
}
