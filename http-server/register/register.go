package register

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stdatiks/jdenticon-go"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/slog"
	"net/http"
	"os"
	"social_media_backend/storage/models"
	"social_media_backend/storage/postgresql"
	"time"
)

func New(log *slog.Logger, storage *postgresql.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User

		if err := c.ShouldBind(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fmt.Println(user.DateOfBirth)

		if user.Email == "" || user.Password == "" || user.Name == "" {
			c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"error": "Email or Password or Name is empty"})
			log.Error("Email or Password or Name is empty")
			return
		}

		isUserExist := storage.IsExistByEmail(user.Email)
		if isUserExist {
			c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"error": " This email Already Exist"})
			log.Error("User with this Email already exists")
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		user.Password = string(hashedPassword)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Error("Error hashing password")
			return
		}

		icon := jdenticon.New(user.Name)
		svg, err := icon.SVG()
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Error("Error generating SVG")
			return
		}
		avatarName := fmt.Sprintf("%s_%s.svg", user.Name, time.Now().Format("2006_01_02_15_04_05"))
		avatarPath := "uploads/" + avatarName
		file, err := os.Create(avatarPath)
		user.AvatarURL = avatarPath

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Error("Error creating avatar")
			return
		}
		defer file.Close()

		svgString := string(svg)
		_, err = file.WriteString(svgString)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Error("Error creating avatar")
			return
		}

		err = storage.DB.Create(&user).Error
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Error("Error when creating user")
			return
		}

		c.IndentedJSON(http.StatusCreated, gin.H{"user": user})
	}
}
