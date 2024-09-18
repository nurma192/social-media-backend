package actions

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"social_media_backend/storage/models"
	"strconv"
)

func ParseToUint(value any) (uint, bool) {
	paramID, err := strconv.ParseUint(value.(string), 10, 64)
	if err != nil {
		return 0, false
	}

	return uint(paramID), true
}

func GetUserFromReq(c *gin.Context) models.User {
	object, _ := c.Get("user")
	currentUser, ok := object.(models.User)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "User not found in request context (actions.getUserFromContext)"})
		return models.User{}
	}
	return currentUser
}
