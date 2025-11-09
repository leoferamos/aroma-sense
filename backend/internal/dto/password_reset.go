package dto

// ResetPasswordRequestRequest represents a request to initiate password reset
type ResetPasswordRequestRequest struct {
	Email string `json:"email" binding:"required,email" example:"user@example.com"`
}

// ResetPasswordConfirmRequest represents a request to confirm password reset with code
type ResetPasswordConfirmRequest struct {
	Email       string `json:"email" binding:"required,email" example:"user@example.com"`
	Code        string `json:"code" binding:"required,len=6" example:"123456"`
	NewPassword string `json:"new_password" binding:"required,min=8" example:"NewPass123"`
}
