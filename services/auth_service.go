package services

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/stdatiks/jdenticon-go"
	"golang.org/x/exp/rand"
	"net/http"
	"os"
	"social-media-back/lib/hashing"
	"social-media-back/models"
	"social-media-back/models/request"
	"social-media-back/models/response"
	"time"
)

func (s *AppService) Login(email, password string) (*response.LoginResponse, int, *response.DefaultErrorResponse) {
	user, err := s.getUserByEmail(email)
	if err != nil {
		return nil, http.StatusInternalServerError, &response.DefaultErrorResponse{
			Message: "Server Error",
			Detail:  err.Error(),
		}
	}

	if user == nil {
		return nil, http.StatusUnauthorized, &response.DefaultErrorResponse{
			Message: "Incorrect Email or Password",
		}
	}

	if !hashing.CheckPassword(user.Password, password) {
		return nil, http.StatusUnauthorized, &response.DefaultErrorResponse{
			Message: "Incorrect Email or Password",
		}
	}

	token, err := s.JWTService.GenerateAccessToken(email)
	if err != nil {
		return nil, http.StatusInternalServerError, &response.DefaultErrorResponse{
			Message: "Server Error when try to generate access token",
			Detail:  err.Error(),
		}
	}

	refreshToken, err := s.JWTService.GenerateRefreshToken(email)
	if err != nil {
		return nil, http.StatusInternalServerError, &response.DefaultErrorResponse{
			Message: "Server Error, when try to generate refresh token",
			Detail:  err.Error(),
		}
	}

	fmt.Println("login refreshToken", refreshToken)

	res := &response.LoginResponse{
		RefreshToken: refreshToken,
		Token:        token,
		Success:      true,
	}
	return res, http.StatusOK, nil

}

func (s *AppService) Register(request request.RegisterRequest) (*response.DefaultSuccessResponse, int, *response.DefaultErrorResponse) {
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

	_, code, errRes := s.SendVerifyCode(request.Email)
	if errRes != nil {
		return nil, code, errRes
	}

	err = s.RedisService.SaveRegisteredUserData(&request)
	if err != nil {
		return nil, http.StatusInternalServerError, &response.DefaultErrorResponse{
			Message: "Server Error",
			Detail:  err.Error(),
		}
	}

	res := &response.DefaultSuccessResponse{
		Message: "Code sent to email, verify your account",
		Success: true,
	}

	return res, http.StatusOK, nil
}

func (s *AppService) SendVerifyCode(email string) (*response.DefaultSuccessResponse, int, *response.DefaultErrorResponse) {
	//todo Save code logic
	codeExist, err := s.RedisService.CheckVerificationCode(email)
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
	err = s.RedisService.SetVerificationCode(email, code)
	if err != nil {
		return nil, http.StatusInternalServerError, &response.DefaultErrorResponse{
			Message: "Server Error, failed to save verification code",
			Detail:  fmt.Sprintf("failed to save verification code: %w", err.Error()),
		}
	}
	//todo Send code to email logic
	err = s.sendMessage(email, "Your verification code", "Your verification code: "+code+"\n Verify Your account within 10 minutes, Registration tickets time out after 10 minutes")
	if err != nil {
		return nil, http.StatusInternalServerError, &response.DefaultErrorResponse{
			Message: "Server Error, failed to send verification code",
			Detail:  fmt.Sprintf("failed to send verification code: %w", err.Error()),
		}
	}

	res := &response.DefaultSuccessResponse{
		Success: true,
		Message: "Verify code successfully sent to your email",
	}
	return res, http.StatusOK, nil
}

func (s *AppService) VerifyAccount(email, code string) (*response.VerifyAccountResponse, int, *response.DefaultErrorResponse) {
	storedCode, err := s.RedisService.GetVerificationCode(email)
	if err != nil {
		return nil, http.StatusInternalServerError, &response.DefaultErrorResponse{
			Message: "Server Error",
			Detail:  fmt.Sprintf("Error when try to get verification code from redisStorage: %w", err.Error()),
		}
	}

	if storedCode == "" {
		return nil, http.StatusNotFound, &response.DefaultErrorResponse{
			Message: "User not found",
			Detail:  fmt.Sprintf("User %s not found ", email),
		}
	}

	if storedCode != code {
		return nil, http.StatusForbidden, &response.DefaultErrorResponse{
			Message: "Wrong verification code",
			Detail:  fmt.Sprintf("wrong verification code"),
		}
	}

	userData, err := s.RedisService.GetRegisteredUserByEmail(email)
	if err != nil {
		return nil, http.StatusInternalServerError, &response.DefaultErrorResponse{
			Message: "Server Error",
			Detail:  "Error when try to get Registered User ticket " + err.Error(),
		}
	}

	if userData == nil {
		return nil, http.StatusNotFound, &response.DefaultErrorResponse{
			Message: "Registration tickets time is out",
			Detail:  fmt.Sprintf("Registration tickets time is out"),
		}
	}

	hashedPassword, err := hashing.HashPassword(userData.Password)
	if err != nil {
		return nil, http.StatusInternalServerError, &response.DefaultErrorResponse{
			Message: "Server Error",
			Detail:  fmt.Sprintf("failed to hashing password: %w", err.Error()),
		}
	}

	icon := jdenticon.New(userData.Firstname)
	svg, err := icon.SVG()
	if err != nil {
		return nil, http.StatusInternalServerError, &response.DefaultErrorResponse{
			Message: "Server Error",
			Detail:  fmt.Sprintf("failed to generate svg: %w", err.Error()),
		}
	}
	avatarName := fmt.Sprintf("%s_%s.svg", userData.Firstname, time.Now().Format("2006_01_02_15_04_05"))
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
		userData.Email,
		hashedPassword,
		userData.Firstname,
		userData.Lastname,
		avatarURL,
	).Scan(&userID)

	if err != nil {
		return nil, http.StatusInternalServerError, &response.DefaultErrorResponse{
			Message: "Server Error",
			Detail:  fmt.Sprintf("failed to insert user into DB: %w", err.Error()),
		}
	}
	createdUser := &models.User{
		ID:        userID,
		Email:     userData.Email,
		Firstname: userData.Firstname,
		Lastname:  userData.Lastname,
		AvatarURL: &avatarURL,
		CreatedAt: time.Now(),
	}

	_, err = s.DB.Exec("UPDATE users SET verified = true WHERE email = $1", email)
	if err != nil {
		return nil, http.StatusInternalServerError, &response.DefaultErrorResponse{
			Message: "Server Error",
			Detail:  fmt.Sprintf("failed to update verification of account: %w", err.Error()),
		}
	}

	_ = s.RedisService.DeleteVerificationCode(email)
	_ = s.RedisService.DeleteRegisteredUserByEmail(email)

	res := &response.VerifyAccountResponse{
		User:    createdUser,
		Success: true,
		Message: "Your account has been verified successfully!",
	}
	return res, http.StatusCreated, nil
}

func (s *AppService) RefreshToken(refreshToken string) (*response.LoginResponse, int, *response.DefaultErrorResponse) {
	claims, err := s.JWTService.ValidateToken(refreshToken)
	fmt.Println("refreshToken:", refreshToken)
	if err != nil {
		return nil, http.StatusInternalServerError, &response.DefaultErrorResponse{
			Message: "Server Error",
			Detail:  fmt.Sprintf("failed to validate token: %w", err.Error()),
		}
	}

	token, err := s.JWTService.GenerateAccessToken(claims.Email)
	if err != nil {
		return nil, http.StatusInternalServerError, &response.DefaultErrorResponse{
			Message: "Server Error",
			Detail:  fmt.Sprintf("failed to generate access token: %w", err.Error()),
		}
	}
	newRefreshToken, err := s.JWTService.GenerateRefreshToken(claims.Email)
	if err != nil {
		return nil, http.StatusInternalServerError, &response.DefaultErrorResponse{
			Message: "Server Error, when try to generate refresh token",
			Detail:  err.Error(),
		}
	}

	res := &response.LoginResponse{
		RefreshToken: newRefreshToken,
		Token:        token,
		Success:      true,
	}
	return res, http.StatusOK, nil
}
