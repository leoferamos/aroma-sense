package dto

import "time"

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

// UpdateProfileRequest represents the payload to update user's profile fields.
type UpdateProfileRequest struct {
	DisplayName string `json:"display_name" binding:"required,min=2,max=50" example:"João Santos"`
}

// AdminDeactivateUserRequest represents the payload for admin user deactivation with enhanced LGPD compliance
type AdminDeactivateUserRequest struct {
	Reason          string     `json:"reason" binding:"required,oneof=violation_of_terms privacy_violation fraud_suspicion account_compromise underage_user duplicate_account" example:"violation_of_terms"`
	Notes           string     `json:"notes" binding:"max=500" example:"User violated terms of service by posting inappropriate content"`
	SuspensionUntil *time.Time `json:"suspension_until,omitempty" example:"2024-12-31T23:59:59Z"`
}

// ContestationRequest represents the payload for user contestation of account deactivation
type ContestationRequest struct {
	Reason string `json:"reason" binding:"required,min=10,max=500" example:"I believe my account was deactivated by mistake. I did not violate any terms of service."`
}

// AdminReactivateUserRequest represents the payload for admin user reactivation
type AdminReactivateUserRequest struct {
	Reason string `json:"reason" binding:"required,min=10,max=200" example:"Contestation reviewed and approved - account reactivated"`
}
