package getUserByID

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
	"net/http"
	"social_media_backend/lib/actions"
	"social_media_backend/storage/models"
	"social_media_backend/storage/postgresql"
)

func New(log *slog.Logger, storage *postgresql.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUser := actions.GetUserFromReq(c)
		if currentUser.ID == 0 {
			return
		}

		paramID := c.Param("id")
		var userFromParam *models.User
		err := storage.DB.
			Preload("Followings").
			Preload("Followers").
			First(&userFromParam, paramID).Error
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}

		if userFromParam == nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
			return
		}

		isIFollowedToHim := storage.IsThisUserFollowedTo(currentUser.ID, userFromParam.ID)
		isHeFollowedToMe := storage.IsThisUserFollowedTo(userFromParam.ID, currentUser.ID)

		c.IndentedJSON(http.StatusOK, gin.H{
			"user":             userFromParam,
			"isIFollowedToHim": isIFollowedToHim,
			"isHeFollowedToMe": isHeFollowedToMe,
		})

	}
}
