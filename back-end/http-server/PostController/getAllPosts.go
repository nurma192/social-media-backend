package postController

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
	"net/http"
	"social_media_backend/lib/actions"
	"social_media_backend/storage/models"
	"social_media_backend/storage/postgresql"
)

type PostWithLikes struct {
	Post          models.Post `json:"post"`
	IsLikedByUser bool        `json:"isLikedByUser"`
}

func GetAllPosts(log *slog.Logger, storage *postgresql.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUser := actions.GetUserFromReq(c)
		if currentUser.ID == 0 {
			return
		}

		var posts []models.Post
		err := storage.DB.Order("created_at desc").Preload("Likes").Preload("Comments").Find(&posts).Error
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Error("Failed to get all posts", "error", err)
			return
		}

		if len(posts) == 0 {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Posts not found"})
			return
		}

		var postWithLikes []PostWithLikes
		for _, post := range posts {
			isLiked := false
			for _, like := range post.Likes {
				if like.UserID == currentUser.ID {
					isLiked = true
					continue
				}
			}
			postWithLikes = append(postWithLikes, PostWithLikes{
				Post:          post,
				IsLikedByUser: isLiked,
			})
		}

		c.IndentedJSON(http.StatusOK, postWithLikes)
	}
}
