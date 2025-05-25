package entity

type ApiResponse struct {
	Success bool          `json:"success"`
	Error   ErrorResponse `json:"error,omitempty"`
	Data    interface{}   `json:"data,omitempty"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Type    string `json:"type"`
	Message string `json:"message"`
}
