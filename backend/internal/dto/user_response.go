package dto

import (
	"time"
)

// UserResponse represents the public data of a user returned by the API.
type UserResponse struct {
	PublicID  string    `json:"public_id" example:"uuid"`
	Email     string    `json:"email" example:"user@example.com"`
	Role      string    `json:"role" example:"user"`
	CreatedAt time.Time `json:"created_at" example:"2025-09-11T12:00:00Z"`
}

// LoginResponse represents a successful login response containing a JWT token and user info.
type LoginResponse struct {
	Message     string       `json:"message" example:"Login successful"`
	AccessToken string       `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User        UserResponse `json:"user"`
}

// ProfileResponse represents the current user's profile data
type ProfileResponse struct {
	PublicID    string    `json:"public_id"`
	Email       string    `json:"email"`
	Role        string    `json:"role"`
	DisplayName *string   `json:"display_name,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// UserExportResponse represents all user data for GDPR export
type UserExportResponse struct {
	PublicID            string     `json:"public_id"`
	Email               string     `json:"email"`
	Role                string     `json:"role"`
	DisplayName         *string    `json:"display_name,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	LastLoginAt         *time.Time `json:"last_login_at,omitempty"`
	DeactivatedAt       *time.Time `json:"deactivated_at,omitempty"`
	DeletionRequestedAt *time.Time `json:"deletion_requested_at,omitempty"`
	DeletionConfirmedAt *time.Time `json:"deletion_confirmed_at,omitempty"`
}

// AdminUserResponse represents user data for admin interface
type AdminUserResponse struct {
	ID                    uint       `json:"id"`
	PublicID              string     `json:"public_id"`
	MaskedEmail           string     `json:"masked_email"`
	Role                  string     `json:"role"`
	DisplayName           *string    `json:"display_name,omitempty"`
	CreatedAt             time.Time  `json:"created_at"`
	LastLoginAt           *time.Time `json:"last_login_at,omitempty"`
	DeactivatedAt         *time.Time `json:"deactivated_at,omitempty"`
	DeactivatedBy         *string    `json:"deactivated_by,omitempty"`
	DeactivationReason    *string    `json:"deactivation_reason,omitempty"`
	DeactivationNotes     *string    `json:"deactivation_notes,omitempty"`
	SuspensionUntil       *time.Time `json:"suspension_until,omitempty"`
	ReactivationRequested bool       `json:"reactivation_requested,omitempty"`
	ContestationDeadline  *time.Time `json:"contestation_deadline,omitempty"`
}

// UserListRequest represents pagination and filter parameters for user listing
type UserListRequest struct {
	Limit   int                    `json:"limit" example:"10"`
	Offset  int                    `json:"offset" example:"0"`
	Filters map[string]interface{} `json:"filters,omitempty"`
}

// UserListResponse represents paginated user list for admin
type UserListResponse struct {
	Users  []AdminUserResponse `json:"users"`
	Total  int64               `json:"total"`
	Limit  int                 `json:"limit"`
	Offset int                 `json:"offset"`
}

// UpdateRoleRequest represents role update request
type UpdateRoleRequest struct {
	Role string `json:"role" example:"admin" validate:"required,oneof=admin client"`
}

// DeleteAccountRequest represents account deletion confirmation
type DeleteAccountRequest struct {
	Confirmation string `json:"confirmation" example:"DELETE_MY_ACCOUNT or EXCLUIR_MINHA_CONTA" validate:"required,oneof=DELETE_MY_ACCOUNT EXCLUIR_MINHA_CONTA"`
}
