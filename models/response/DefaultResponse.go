package response

type DefaultResponse struct {
	Message string `json:"message"`
	Detail  string `json:"detail"`
	Success bool   `json:"success" `
}

type Response struct {
	Result  interface{} `json:"result,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type PaginationResponse struct {
	Page       int         `json:"page,omitempty"`
	TotalPages int         `json:"totalPages,omitempty"`
	Limit      int         `json:"limit,omitempty"`
	Result     interface{} `json:"result,omitempty"`
}
