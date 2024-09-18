package followController

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
	"net/http"
	"social_media_backend/lib/actions"
	"social_media_backend/storage/models"
	"social_media_backend/storage/postgresql"
)

type Follow struct {
	UserID     uint
	FollowerID uint
}

type FollowReqStruct struct {
	FollowingID uint
}

func FollowUser(log *slog.Logger, storage *postgresql.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var followReqStruct FollowReqStruct
		if err := c.ShouldBind(&followReqStruct); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		currentUser := actions.GetUserFromReq(c)
		if currentUser.ID == 0 {
			return
		}

		if currentUser.ID == followReqStruct.FollowingID {
			c.IndentedJSON(http.StatusConflict, gin.H{"error": "You can't follow yourself"})
		}

		var count int64
		if err := storage.DB.Limit(1).Table("follows").Where("follower_id = ? AND user_id = ?", currentUser.ID, followReqStruct.FollowingID).Count(&count).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			log.Error("Error when try to find is follow exist: ", err)
			return
		}

		if count > 0 {
			c.IndentedJSON(http.StatusConflict, gin.H{"error": "You already followed to this user"})
			return
		}

		if err := storage.DB.Model(&currentUser).Association("Followings").Append(&models.User{ID: followReqStruct.FollowingID}); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			log.Error("Error when try to add following user: ", err)
			return
		}

		c.IndentedJSON(http.StatusCreated, gin.H{
			"message": "User successfully followed to user",
		})

	}
}
