package dto

// CreateUserRequest represents the expected payload for registering a new user.
type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com" format:"email"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginRequest represents the expected payload for user login.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com" format:"email"`
	Password string `json:"password" binding:"required"`
}
