package login

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/exp/slog"
	"net/http"
	"os"
	"social_media_backend/lib/hash"
	"social_media_backend/storage/models"
	"social_media_backend/storage/postgresql"
	"time"
)

type login struct {
	Email    string
	Password string
}

func New(log *slog.Logger, storage *postgresql.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userLogin login

		if err := c.ShouldBind(&userLogin); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if userLogin.Email == "" || userLogin.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email or password is empty"})
			return
		}

		isUserExist := storage.IsExistByEmail(userLogin.Email)

		if !isUserExist {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User does not exist"})
			return
		}

		user, err := storage.GetUserBy(models.User{Email: userLogin.Email})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		isPasswordEqual := hash.CheckPasswordHash(userLogin.Password, user.Password)
		if !isPasswordEqual {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email or password is wrong"})
			return
		}

		// Jwt token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": user.ID,
			"exp": time.Now().Add(time.Hour * 24).Unix(),
		})

		// Sign and get the complete encoded token as a string using the secret
		secret := os.Getenv("SECRET")
		fmt.Println("secret:", secret)
		tokenString, err := token.SignedString([]byte(secret))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("Authorization", tokenString, 3600*24, "", "", false, true)
		c.IndentedJSON(http.StatusOK, gin.H{
			"message": "Login successful",
			"user":    user,
		})
	}
}
