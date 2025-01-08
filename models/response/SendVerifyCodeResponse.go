package response

type SendVerifyCodeResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
