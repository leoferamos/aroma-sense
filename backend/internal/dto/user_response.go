package dto

// UserResponse represents the public data of a user returned by the API.
type UserResponse struct {
	PublicID  string `json:"public_id" example:"uuid"`
	Email     string `json:"email" example:"user@example.com"`
	Role      string `json:"role" example:"user"`
	CreatedAt string `json:"created_at" example:"2025-09-11T12:00:00Z"`
}
