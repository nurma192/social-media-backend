package followController

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
	"net/http"
	"social_media_backend/lib/actions"
	"social_media_backend/storage/models"
	"social_media_backend/storage/postgresql"
)

func UnFollowUser(log *slog.Logger, storage *postgresql.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		followingID, ok := actions.ParseToUint(c.Param("id"))
		if !ok {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		currentUser := actions.GetUserFromReq(c)
		if currentUser.ID == 0 {
			return
		}

		if currentUser.ID == followingID {
			c.IndentedJSON(http.StatusConflict, gin.H{"error": "You can't unFollow from yourself"})
			return
		}

		var count int64
		if err := storage.DB.Limit(1).Table("follows").Where("follower_id = ? AND user_id = ?", currentUser.ID, followingID).Count(&count).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			log.Error("Error when try to find is follow exist: ", err)
			return
		}

		if count == 0 {
			c.IndentedJSON(http.StatusConflict, gin.H{"error": "Сен брат даже подписка жасамағансынғой, қалай отписка жасайсын айтшы :("})
			return
		}

		if err := storage.DB.Model(&currentUser).Association("Followings").Delete(&models.User{ID: followingID}); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			log.Error("Error when try to add following user: ", err)
			return
		}

		c.IndentedJSON(http.StatusCreated, gin.H{
			"message": "User successfully unFollowed from user",
		})
	}
}
