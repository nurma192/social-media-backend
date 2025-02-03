package request

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username  string `json:"username" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Firstname string `json:"firstname" binding:"required"`
	Lastname  string `json:"lastname" binding:"required"`
	Password  string `json:"password" binding:"required,min=1"`
}

type SendVerifyCodeRequest struct {
	Email string `json:"email" binding:"required" binding:"required"`
}

type VerifyAccountRequest struct {
	Email string `json:"email" binding:"required" binding:"required"`
	Code  string `json:"code" binding:"required" binding:"required"`
}
