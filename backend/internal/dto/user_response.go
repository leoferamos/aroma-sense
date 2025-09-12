package dto

import "time"

// UserResponse represents the public data of a user returned by the API.
type UserResponse struct {
	PublicID  string    `json:"public_id" example:"uuid"`
	Email     string    `json:"email" example:"user@example.com"`
	Role      string    `json:"role" example:"user"`
	CreatedAt time.Time `json:"created_at" example:"2025-09-11T12:00:00Z"`
}

// LoginResponse represents a successful login response containing a JWT token and user info.
type LoginResponse struct {
	Token string       `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  UserResponse `json:"user"`
}
