package request

type SendVerifyCodeRequest struct {
	Email string `json:"email" binding:"required"`
}
