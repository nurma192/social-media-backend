package response

type ErrorResponse struct {
	Message string `json:"message"`
	Detail  string `json:"detail"`
	Success bool   `json:"success" `
}
