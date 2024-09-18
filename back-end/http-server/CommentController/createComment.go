package commentController

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
	"net/http"
	"social_media_backend/lib/actions"
	"social_media_backend/storage/models"
	"social_media_backend/storage/postgresql"
)

type CommentReqStruct struct {
	Content string
	PostID  uint
}

func CreateComment(log *slog.Logger, storage *postgresql.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var commentRequest CommentReqStruct
		if err := c.ShouldBind(&commentRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if commentRequest.Content == "" || commentRequest.PostID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "content, post id is empty"})
			return
		}

		currentUser := actions.GetUserFromReq(c)
		if currentUser.ID == 0 {
			return
		}

		comment := models.Comment{
			UserID:  currentUser.ID,
			PostID:  commentRequest.PostID,
			Content: commentRequest.Content,
		}

		err := storage.DB.Create(&comment).Error
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			log.Error("Failed to create comment", err)
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{
			"Comment": comment,
			"message": "Comment created successfully",
		})
	}
}
