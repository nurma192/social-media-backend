package current

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
	"net/http"
	"social_media_backend/lib/actions"
	"social_media_backend/storage/postgresql"
)

func New(log *slog.Logger, storage *postgresql.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUser := actions.GetUserFromReq(c)
		if currentUser.ID == 0 {
			return
		}

		err := storage.DB.
			Preload("Followers").
			Preload("Followings").
			Preload("Posts").
			First(&currentUser).Error

		if err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.IndentedJSON(http.StatusOK, currentUser)

	}
}
