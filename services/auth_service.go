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
	"social-media-back/models/response"
	"time"
)

//message, code, loginResponse, err
// (string, int, *response.LoginResponse, error

func (s *AppService) Login(username, password string) string {
	return "Hello login"
}

func (s *AppService) Register(request request.RegisterRequest) (*response.RegisterResponse, int, *response.DefaultErrorResponse) {
	var existingUserId int
	err := s.DB.QueryRow("SELECT id FROM users where email = $1", request.Email).Scan(&existingUserId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, http.StatusInternalServerError, &response.DefaultErrorResponse{
			Message: "Server Error",
			Detail:  fmt.Sprintf("failed to check existing user: %w", err.Error()),
		}
	}

	if err == nil {
		return nil, http.StatusConflict, &response.DefaultErrorResponse{
			Message: "User with this email already exists",
			Detail:  "Conflict",
		}
	}

	hashedPassword, err := hash.HashPassword(request.Password)
	if err != nil {
		return nil, http.StatusInternalServerError, &response.DefaultErrorResponse{
			Message: "Server Error",
			Detail:  fmt.Sprintf("failed to hash password: %w", err.Error()),
		}
	}

	icon := jdenticon.New(request.Firstname)
	svg, err := icon.SVG()
	if err != nil {
		return nil, http.StatusInternalServerError, &response.DefaultErrorResponse{
			Message: "Server Error",
			Detail:  fmt.Sprintf("failed to generate svg: %w", err.Error()),
		}
	}
	avatarName := fmt.Sprintf("%s_%s.svg", request.Firstname, time.Now().Format("2006_01_02_15_04_05"))
	avatarPath := "uploads/avatars/" + avatarName
	file, err := os.Create(avatarPath)

	if err != nil {
		return nil, http.StatusInternalServerError, &response.DefaultErrorResponse{
			Message: "Server Error",
			Detail:  fmt.Sprintf("failed to create avatar file: %w", err.Error()),
		}
	}
	defer file.Close()

	svgString := string(svg)
	_, err = file.WriteString(svgString)
	if err != nil {
		return nil, http.StatusInternalServerError, &response.DefaultErrorResponse{
			Message: "Server Error",
			Detail:  fmt.Sprintf("failed to write avatar file: %w", err.Error()),
		}
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
		return nil, http.StatusInternalServerError, &response.DefaultErrorResponse{
			Message: "Server Error",
			Detail:  fmt.Sprintf("failed to insert user into DB: %w", err.Error()),
		}
	}
	user := &models.User{
		ID:        userID,
		Email:     request.Email,
		Firstname: request.Firstname,
		Lastname:  request.Lastname,
		AvatarURL: &avatarURL,
		CreatedAt: time.Now(),
	}

	res := &response.RegisterResponse{
		User:    user,
		Message: "User successfully created",
		Success: true,
	}

	return res, http.StatusCreated, nil
}

// will return message, statusCode, err
func (s *AppService) SendVerifyCode(email string) (*response.SendVerifyCodeResponse, int, *response.DefaultErrorResponse) {
	user, err := s.getUserByEmail(email)
	if err != nil {
		return nil, http.StatusInternalServerError, &response.DefaultErrorResponse{
			Message: "Server error",
			Detail:  err.Error(),
		}
	}
	if user == nil {
		return nil, http.StatusNotFound, &response.DefaultErrorResponse{
			Message: "User not found",
			Detail:  fmt.Sprintf("User not found: %s", email),
		}
	}
	if user.Verified {
		return nil, http.StatusForbidden, &response.DefaultErrorResponse{
			Message: "User is already verified",
			Detail:  fmt.Sprintf("User is already verified: %s", email),
		}
	}

	//todo Save code logic
	codeExist, err := s.checkVerificationCode(email)
	if err != nil {
		return nil, http.StatusInternalServerError, &response.DefaultErrorResponse{
			Message: "Server Error",
			Detail:  fmt.Sprintf("failed to check verification code: %w", err.Error()),
		}
	}
	if codeExist {
		return nil, http.StatusConflict, &response.DefaultErrorResponse{
			Message: "User already verified",
			Detail:  fmt.Sprintf("User already verified: %s", email),
		}
	}

	rand.Seed(uint64(time.Now().UnixNano()))
	code := fmt.Sprintf("%04d", rand.Intn(10000))
	err = s.setVerificationCode(email, code)
	if err != nil {
		return nil, http.StatusInternalServerError, &response.DefaultErrorResponse{
			Message: "Server Error",
			Detail:  fmt.Sprintf("failed to set verification code: %w", err.Error()),
		}
	}
	//todo Send code to email logic
	s.sendMessage(email)

	res := &response.SendVerifyCodeResponse{
		Success: true,
		Message: "Verify code successfully sent to your email: " + code,
	}
	return res, http.StatusOK, nil
}

func (s *AppService) VerifyAccount(email, code string) (*response.VerifyAccountResponse, int, *response.DefaultErrorResponse) {
	storedCode, err := s.getVerificationCode(email)
	if err != nil {
		return nil, http.StatusInternalServerError, &response.DefaultErrorResponse{
			Message: "Server Error",
			Detail:  fmt.Sprintf("Error when try to get verification code from redis: %w", err.Error()),
		}
	}

	if storedCode != code {
		return nil, http.StatusForbidden, &response.DefaultErrorResponse{
			Message: "Wrong verification code",
			Detail:  fmt.Sprintf("wrong verification code"),
		}
	}

	_, err = s.DB.Exec("UPDATE users SET verified = true WHERE email = $1", email)
	if err != nil {
		return nil, http.StatusInternalServerError, &response.DefaultErrorResponse{
			Message: "Server Error",
			Detail:  fmt.Sprintf("failed to update verification of account: %w", err.Error()),
		}
	}

	err = s.deleteVerificationCode(email)

	res := &response.VerifyAccountResponse{
		Success: true,
		Message: "Your account has been verified successfully!",
	}
	return res, http.StatusOK, nil
}
