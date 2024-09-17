package postController

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
	"net/http"
	"social_media_backend/storage/models"
	"social_media_backend/storage/postgresql"
)

type PostReqStruct struct {
	Content string `json:"content"`
}

func CreatePost(log *slog.Logger, storage *postgresql.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var postReq PostReqStruct
		if err := c.ShouldBind(&postReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if postReq.Content == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "content is required"})
			return
		}

		object, _ := c.Get("user")

		currentUser, ok := object.(models.User)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "User not found in request context"})
			log.Error("User not found in request context")
			return
		}
		newPost := models.Post{
			Content: postReq.Content,
			UserID:  currentUser.ID,
		}
		err := storage.DB.Create(&newPost).Error
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Error("Failed to create post", "error", err)
			return
		}

		c.IndentedJSON(http.StatusOK, newPost)
	}
}
