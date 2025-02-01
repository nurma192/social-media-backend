package services

import (
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

func (s *AppService) Login(email, password string) (*response.Response, int, string) {
	user, err := s.DBService.GetUserByEmail(email)
	if err != nil {
		return &response.Response{
			Error: err.Error(),
		}, http.StatusInternalServerError, ""
	}

	if user == nil {
		return &response.Response{
			Error: "Incorrect Email or Password",
		}, http.StatusBadRequest, ""
	}

	if !hashing.CheckPassword(user.Password, password) {
		return &response.Response{
			Error: "Incorrect Email or Password",
		}, http.StatusBadRequest, ""
	}

	token, err := s.JWTService.GenerateAccessToken(email, user.Id)
	if err != nil {
		return &response.Response{
			Error: "Server Error when try to generate access token",
		}, http.StatusInternalServerError, ""
	}

	refreshToken, err := s.JWTService.GenerateRefreshToken(email)
	if err != nil {
		return &response.Response{
			Error: err.Error(),
		}, http.StatusInternalServerError, ""
	}

	fmt.Println("login refreshToken", refreshToken)

	loginResponse := &response.LoginResponse{
		Token: token,
	}
	return &response.Response{
		Result:  loginResponse,
		Message: "Login Success",
	}, http.StatusOK, refreshToken

}

func (s *AppService) Register(request request.RegisterRequest) (*response.Response, int) {
	isExistByEmail, err := s.DBService.IsUserExistByEmail(request.Email)
	isExistByUsername, err := s.DBService.IsUserExistByUsername(request.Username)

	if isExistByEmail {
		return &response.Response{
			Error: "User With this Email Already Exist",
		}, http.StatusBadRequest
	}
	if isExistByUsername {
		return &response.Response{
			Error: "User With this Username Already Exist",
		}, http.StatusBadRequest
	}

	err = s.RedisService.SaveRegisteredUserData(&request)
	if err != nil {
		return &response.Response{
			Error: err.Error(),
		}, http.StatusInternalServerError
	}

	res, code := s.SendVerifyCode(request.Email)
	fmt.Println(res)
	if res.Error != "" {
		return res, code
	}

	return &response.Response{
		Message: "Code sent to email, verify your account",
	}, http.StatusOK
}

func (s *AppService) SendVerifyCode(email string) (*response.Response, int) {
	exist, err := s.DBService.IsUserExistByEmail(email)
	if err != nil {
		return &response.Response{
			Error: err.Error(),
		}, http.StatusInternalServerError
	}
	if exist {
		return &response.Response{
			Error: "User With this Email Already Exist",
		}, http.StatusBadRequest
	}
	userData, err := s.RedisService.GetRegisteredUserByEmail(email)
	if err != nil {
		return &response.Response{
			Error: "Cant check the users register ticket" + err.Error(),
		}, http.StatusInternalServerError
	}
	if userData == nil {
		return &response.Response{
			Error: "Register ticket not found",
		}, http.StatusBadRequest
	}
	//todo Save code logic
	codeExist, err := s.RedisService.CheckVerificationCode(email)
	if err != nil {
		return &response.Response{
			Error: fmt.Sprintf("failed to check verification code: %w", err.Error()),
		}, http.StatusInternalServerError
	}
	if codeExist {
		return &response.Response{
			Error: fmt.Sprintf("User verify code already sent: %s", email),
		}, http.StatusConflict
	}

	rand.Seed(uint64(time.Now().UnixNano()))
	code := fmt.Sprintf("%04d", rand.Intn(10000))
	err = s.RedisService.SetVerificationCode(email, code)
	if err != nil {
		return &response.Response{
			Error: fmt.Sprintf("failed to save verification code: %w", err.Error()),
		}, http.StatusInternalServerError
	}
	//todo Send code to email logic
	err = s.EmailService.SendMessage(email, "Your verification code", "Your verification code: "+code+"\n Verify Your account within 10 minutes, Registration tickets time out after 10 minutes")
	if err != nil {
		return &response.Response{
			Error: fmt.Sprintf("failed to send verification code: %w", err.Error()),
		}, http.StatusInternalServerError
	}

	return &response.Response{
		Message: "Verify code successfully sent to your email",
	}, http.StatusOK
}

func (s *AppService) VerifyAccount(email, code string) (*response.Response, int) {
	storedCode, err := s.RedisService.GetVerificationCode(email)
	if err != nil {
		return &response.Response{
			Error: fmt.Sprintf("Error when try to get verification code from redisStorage: %w", err.Error()),
		}, http.StatusInternalServerError
	}

	if storedCode == "" {
		return &response.Response{
			Error: fmt.Sprintf("User %s not found ", email),
		}, http.StatusNotFound
	}

	if storedCode != code {
		return &response.Response{
			Error: "wrong verification code",
		}, http.StatusForbidden
	}

	userData, err := s.RedisService.GetRegisteredUserByEmail(email)
	if err != nil {
		return &response.Response{
			Error: "Error when try to get Registered User ticket " + err.Error(),
		}, http.StatusInternalServerError
	}

	if userData == nil {
		return &response.Response{
			Error: fmt.Sprintf("Registration tickets time is out"),
		}, http.StatusNotFound
	}

	hashedPassword, err := hashing.HashPassword(userData.Password)
	if err != nil {
		return &response.Response{
			Error: fmt.Sprintf("failed to hashing password: %w", err.Error()),
		}, http.StatusInternalServerError
	}

	icon := jdenticon.New(userData.Firstname)
	svg, err := icon.SVG()
	if err != nil {
		return &response.Response{
			Error: fmt.Sprintf("failed to generate svg: %w", err.Error()),
		}, http.StatusInternalServerError
	}
	avatarName := fmt.Sprintf("%s_%s.svg", userData.Firstname, time.Now().Format("2006_01_02_15_04_05"))
	avatarPath := "uploads/avatars/" + avatarName
	file, err := os.Create(avatarPath)

	if err != nil {
		return &response.Response{
			Error: fmt.Sprintf("failed to create avatar file: %w", err.Error()),
		}, http.StatusInternalServerError
	}
	defer file.Close()

	svgString := string(svg)
	_, err = file.WriteString(svgString)
	if err != nil {
		return &response.Response{
			Error: fmt.Sprintf("failed to write avatar file: %w", err.Error()),
		}, http.StatusInternalServerError
	}

	avatarURL := avatarPath

	query := `INSERT INTO users (email, username, password, firstname, lastname, avatar_url) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	var userId int
	err = s.DBService.DB.QueryRow(
		query,
		userData.Email,
		userData.Username,
		hashedPassword,
		userData.Firstname,
		userData.Lastname,
		avatarURL,
	).Scan(&userId)

	if err != nil {
		return &response.Response{
			Error: fmt.Sprintf("failed to insert user into DB: %w", err.Error()),
		}, http.StatusInternalServerError
	}
	createdUser := &models.User{
		Id:        userId,
		Email:     userData.Email,
		Username:  userData.Username,
		Firstname: userData.Firstname,
		Lastname:  userData.Lastname,
		AvatarURL: &avatarURL,
		CreatedAt: time.Now(),
	}

	_, err = s.DBService.DB.Exec("UPDATE users SET verified = true WHERE email = $1", email)
	if err != nil {
		return &response.Response{
			Error: fmt.Sprintf("failed to update verification of account: %w", err.Error()),
		}, http.StatusInternalServerError
	}

	_ = s.RedisService.DeleteVerificationCode(email)
	_ = s.RedisService.DeleteRegisteredUserByEmail(email)

	verifyAccountRes := response.VerifyAccountResponse{
		User: createdUser,
	}

	return &response.Response{
		Result:  verifyAccountRes,
		Message: "Your account has been verified successfully!",
	}, http.StatusCreated
}

func (s *AppService) RefreshToken(refreshToken string) (*response.Response, int, string) {
	claims, err := s.JWTService.ValidateToken(refreshToken)
	fmt.Println("refreshToken:", refreshToken)
	if err != nil {
		return &response.Response{
			Error: fmt.Sprintf("failed to validate token: %w", err.Error()),
		}, http.StatusInternalServerError, ""
	}

	token, err := s.JWTService.GenerateAccessToken(claims.Email, claims.UserId)
	if err != nil {
		return &response.Response{
			Error: fmt.Sprintf("failed to generate access token: %w", err.Error()),
		}, http.StatusInternalServerError, ""
	}
	newRefreshToken, err := s.JWTService.GenerateRefreshToken(claims.Email)
	if err != nil {
		return &response.Response{
			Error: err.Error(),
		}, http.StatusInternalServerError, ""
	}

	loginRes := &response.LoginResponse{
		Token: token,
	}
	return &response.Response{
		Result: loginRes,
	}, http.StatusOK, newRefreshToken
}
