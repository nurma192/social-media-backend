package current

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
	"net/http"
	"social_media_backend/storage/models"
	"social_media_backend/storage/postgresql"
)

func New(log *slog.Logger, storage *postgresql.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		object, _ := c.Get("user")

		currentUser, ok := object.(models.User)
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		err := storage.DB.
			Preload("Followers").
			Preload("Followings").
			First(&currentUser).Error

		if err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}

		c.IndentedJSON(http.StatusOK, currentUser)

	}
}
