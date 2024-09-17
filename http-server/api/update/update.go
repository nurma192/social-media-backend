package update

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
	"net/http"
	"social_media_backend/lib/actions"
	"social_media_backend/lib/hash"
	"social_media_backend/storage/models"
	"social_media_backend/storage/postgresql"
)

func New(log *slog.Logger, storage *postgresql.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		object, _ := c.Get("user")

		currentUser, ok := object.(models.User)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "User not found in request context"})
			log.Error("User not found in request context")
			return
		}

		paramID, ok := actions.ParseToUint(c.Param("id"))
		if ok == false {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Invalid id"})
			return
		}

		if currentUser.ID != paramID {
			c.IndentedJSON(http.StatusNotFound, gin.H{
				"message": "you dont have permission to update",
			})
			return
		}

		var updatedUser models.User
		if err := c.ShouldBind(&updatedUser); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// todo check email, password, name validation

		if updatedUser.Email != "" {
			if storage.IsExistByEmail(updatedUser.Email) == true {
				c.IndentedJSON(http.StatusConflict, gin.H{
					"message": "user with this email already exists",
				})
				return
			}
		}

		if updatedUser.Password != "" {
			hashedPass, err := hash.HashPassword(updatedUser.Password)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				log.Error("Error when try to hash the password")
				return
			}
			updatedUser.Password = hashedPass
		}

		finallyUser := currentUser
		if updatedUser.Email != "" {
			finallyUser.Email = updatedUser.Email
		}
		if updatedUser.Name != "" {
			finallyUser.Name = updatedUser.Name
		}
		if updatedUser.Password != "" {
			finallyUser.Password = updatedUser.Password
		}
		if updatedUser.Bio != "" {
			finallyUser.Bio = updatedUser.Bio
		}
		if updatedUser.Location != "" {
			finallyUser.Location = updatedUser.Location
		}

		//todo get and set the photo

		if err := storage.DB.Updates(&finallyUser).Error; err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Error("Error when try to update the user!")
			return
		}
		c.IndentedJSON(http.StatusOK, gin.H{
			"status": "ok",
			"user":   finallyUser,
		})
	}
}
