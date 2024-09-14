package getUserByID

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
	"net/http"
	"social_media_backend/storage/models"
	"social_media_backend/storage/postgresql"
)

func New(log *slog.Logger, storage *postgresql.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var user *models.User
		err := storage.DB.First(&user, id).Error
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}

		if user == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "User not found",
			})
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{
			"user": user,
		})

	}
}
