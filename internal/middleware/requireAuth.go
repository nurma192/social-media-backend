package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/exp/slog"
	"net/http"
	"os"
	"social_media_backend/storage/models"
	"social_media_backend/storage/postgresql"
	"time"
)

func RequireAuth(log *slog.Logger, storage *postgresql.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get cookie of rec
		tokenString, err := c.Cookie("Authorization")
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Decode/validate it
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return []byte(os.Getenv("SECRET")), nil
		})
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// check the exp
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			id, ok := claims["sub"].(float64)
			if !ok {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			user, err := storage.GetUserBy(models.User{ID: uint(id)})
			if err != nil || user.ID == 0 {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			// Attach the req
			c.Set("user", user)

			// Continue
			c.Next()
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

	}
}
