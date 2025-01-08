package services

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/stdatiks/jdenticon-go"
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
		return "Server Error", nil, http.StatusInternalServerError, fmt.Errorf("Failed to check existing user: %w", err)
	}

	if err == nil {
		return "User with this email already exists", nil, http.StatusConflict, errors.New("Conflict")
	}

	hashedPassword, err := hash.HashPassword(request.Password)
	if err != nil {
		return "Server Error", nil, http.StatusInternalServerError, fmt.Errorf("Failed tohash the password: %w", err)
	}

	icon := jdenticon.New(request.Firstname)
	svg, err := icon.SVG()
	if err != nil {
		return "Server Error", nil, http.StatusInternalServerError, fmt.Errorf("Failed to generate icon: %w", err)
	}
	avatarName := fmt.Sprintf("%s_%s.svg", request.Firstname, time.Now().Format("2006_01_02_15_04_05"))
	avatarPath := "uploads/avatars/" + avatarName
	file, err := os.Create(avatarPath)

	if err != nil {
		return "Server Error", nil, http.StatusInternalServerError, fmt.Errorf("Failed to create avatar: %w", err)
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

	//todo Send code logic

	return "Verify code successfully sent to your email", http.StatusOK, nil
}

func (s *AppService) VerifyAccount(email, code string) string {
	return fmt.Sprintf("Verify email %s with code %s", email, code)
}
