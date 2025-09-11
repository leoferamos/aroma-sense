package dto

// MessageResponse represents a generic success message.
type MessageResponse struct {
	Message string `json:"message" example:"successfull message"`
}

// ErrorResponse represents an error response.
type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}

// LoginResponse represents a successful login response containing a JWT token and user info.
type LoginResponse struct {
	Token string       `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  UserResponse `json:"user"`
}
