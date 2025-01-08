package response

type DefaultErrorResponse struct {
	Message string `json:"message"`
	Detail  string `json:"detail"`
	Success bool   `json:"success" `
}

type DefaultSuccessResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}
