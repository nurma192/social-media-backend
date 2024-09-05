package login

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
	"social_media_backend/storage/postgresql"
)

func New(log *slog.Logger, storage *postgresql.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
