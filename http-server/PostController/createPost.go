package postController

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
	"net/http"
	"social_media_backend/lib/actions"
	"social_media_backend/storage/models"
	"social_media_backend/storage/postgresql"
)

type PostReqStruct struct {
	Content string
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

		currentUser := actions.GetUserFromReq(c)
		if currentUser.ID == 0 {
			return
		}

		newPost := models.Post{
			Content: postReq.Content,
			UserID:  currentUser.ID,
		}

		if err := storage.DB.Create(&newPost).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Error("Failed to create post", "error", err)
			return
		}

		c.IndentedJSON(http.StatusOK, newPost)
	}
}
