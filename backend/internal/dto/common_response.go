package dto

// MessageResponse represents a generic success message.
type MessageResponse struct {
	Message string `json:"message" example:"success message"`
}

// ErrorResponse represents an error response.
type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}
