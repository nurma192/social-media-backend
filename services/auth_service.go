package services

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/stdatiks/jdenticon-go"
	"golang.org/x/exp/rand"
	"net/http"
	"os"
	"social-media-back/lib/hash"
	"social-media-back/models"
	"social-media-back/models/request"
	"time"
)

//message, user, status, error,

func (s *AppService) Login(username, password string) string {
	return fmt.Sprintf("Hello, %s! This is the Login service.", username)
}

func (s *AppService) Register(request request.RegisterRequest) (string, *models.User, int, error) {
	var existingUserId int
	err := s.DB.QueryRow("SELECT id FROM users where email = $1", request.Email).Scan(&existingUserId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return "Server Error", nil, http.StatusInternalServerError, fmt.Errorf("failed to check existing user: %w", err)
	}

	if err == nil {
		return "User with this email already exists", nil, http.StatusConflict, errors.New("Conflict")
	}

	hashedPassword, err := hash.HashPassword(request.Password)
	if err != nil {
		return "Server Error", nil, http.StatusInternalServerError, fmt.Errorf("failed tohash the password: %w", err)
	}

	icon := jdenticon.New(request.Firstname)
	svg, err := icon.SVG()
	if err != nil {
		return "Server Error", nil, http.StatusInternalServerError, fmt.Errorf("failed to generate icon: %w", err)
	}
	avatarName := fmt.Sprintf("%s_%s.svg", request.Firstname, time.Now().Format("2006_01_02_15_04_05"))
	avatarPath := "uploads/avatars/" + avatarName
	file, err := os.Create(avatarPath)

	if err != nil {
		return "Server Error", nil, http.StatusInternalServerError, fmt.Errorf("failed to create avatar: %w", err)
	}
	defer file.Close()

	svgString := string(svg)
	_, err = file.WriteString(svgString)
	if err != nil {
		return "Server Error", nil, http.StatusInternalServerError, fmt.Errorf("Failed to write avatar: %w", err)
	}

	avatarURL := avatarPath

	query := `INSERT INTO users (email, password, firstname, lastname, avatar_url) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var userID int
	err = s.DB.QueryRow(
		query,
		request.Email,
		hashedPassword,
		request.Firstname,
		request.Lastname,
		avatarURL,
	).Scan(&userID)

	if err != nil {
		return "Server Error", nil, http.StatusInternalServerError, fmt.Errorf("Failed to create user: %w", err)
	}
	user := &models.User{
		ID:        userID,
		Email:     request.Email,
		Firstname: request.Firstname,
		Lastname:  request.Lastname,
		AvatarURL: &avatarURL,
		CreatedAt: time.Now(),
	}

	return "User successfully created", user, http.StatusCreated, nil
}

// will return message, statusCode, err
func (s *AppService) SendVerifyCode(email string) (string, int, error) {
	user, err := s.getUserByEmail(email)
	if err != nil {
		return "Server error", http.StatusInternalServerError, err
	}
	if user == nil {
		return "User not found", http.StatusNotFound, errors.New("User not found")
	}
	if user.Verified {
		return "User already verified", http.StatusForbidden, errors.New("User already verified")
	}

	//todo Save code logic
	codeExist, err := s.checkVerificationCode(email)
	if err != nil {
		return "Server error", http.StatusInternalServerError, err
	}
	if codeExist {
		return "Code already sent to your email", http.StatusForbidden, errors.New("forbidden, You can send code again after 2 minutes")
	}

	rand.Seed(uint64(time.Now().UnixNano()))
	code := fmt.Sprintf("%04d", rand.Intn(10000))
	err = s.setVerificationCode(email, code)
	if err != nil {
		return "Server error, when try to save verification code to redis", http.StatusInternalServerError, err
	}
	//todo Send code to email logic

	return "Verify code successfully sent to your email: " + code, http.StatusOK, nil
}

func (s *AppService) VerifyAccount(email, code string) (string, int, error) {
	storedCode, err := s.getVerificationCode(email)
	if err != nil {
		return "Error when try to get verification code from redis", http.StatusInternalServerError, err
	}

	if storedCode != code {
		fmt.Println("codes:", storedCode, code)
		return "Wrong verification code", http.StatusForbidden, errors.New("wrong verification code")
	}

	_, err = s.DB.Exec("UPDATE users SET verified = true WHERE email = $1", email)
	if err != nil {
		return "Server Error, when try to update verification of account", http.StatusInternalServerError, err
	}

	err = s.deleteVerificationCode(email)
	if err != nil {
		return "Server error while deleting verification code", http.StatusInternalServerError, err
	}

	return "Account verified successfully", http.StatusOK, nil
}
