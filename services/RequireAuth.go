package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"social-media-back/models/response"
)

func (s *AppService) RequireAuth(ctx *gin.Context) {
	authToken := ctx.GetHeader("Authorization")
	fmt.Println("authToken: ", authToken)

	if authToken == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.DefaultErrorResponse{
			Message: "Unauthorized",
			Detail:  "No Access Token provided",
		})
		return
	}

	if len(authToken) > 6 && authToken[:7] == "Bearer " {
		authToken = authToken[7:]
	}

	claims, err := s.JWTService.ValidateToken(authToken)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.DefaultErrorResponse{
			Message: "Unauthorized",
			Detail:  "Invalid or expired token: " + err.Error(),
		})
		return
	}

	ctx.Set("email", claims.Email)

	ctx.Next()
}
