package response

type DefaultResponse struct {
	Message string `json:"message"`
	Detail  string `json:"detail"`
	Success bool   `json:"success" `
}
