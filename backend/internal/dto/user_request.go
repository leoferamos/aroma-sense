package dto

// CreateUserRequest represents the expected payload for registering a new user.
type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest represents the expected payload for user login.
type LoginRequest struct {
	EmailOrUsername string `json:"email_or_username" binding:"required"`
	Password        string `json:"password" binding:"required"`
}
