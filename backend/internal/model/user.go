package model

import (
	"time"

	"gorm.io/gorm"
)

// User represents the database entity for a registered user.
type User struct {
	ID                    uint       `gorm:"primaryKey" json:"-"`
	PublicID              string     `gorm:"type:uuid;not null;uniqueIndex;default:gen_random_uuid()" json:"public_id"`
	Email                 string     `gorm:"size:128;not null;unique" json:"email"`
	PasswordHash          string     `gorm:"size:256;not null" json:"-"`
	Role                  string     `gorm:"size:16;not null;default:client" json:"role"`
	DisplayName           *string    `gorm:"size:64" json:"display_name,omitempty"`
	RefreshTokenHash      *string    `gorm:"size:255" json:"-"`
	RefreshTokenExpiresAt *time.Time `json:"-"`
	CreatedAt             time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt             time.Time  `gorm:"autoUpdateTime" json:"-"`
	DeletedAt             gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete
	LastLoginAt           *time.Time `json:"last_login_at,omitempty"`
	DeactivatedBy         *string    `gorm:"type:uuid" json:"deactivated_by,omitempty"`
	DeactivatedAt         *time.Time `json:"deactivated_at,omitempty"`
	DeactivationReason    *string    `gorm:"size:50" json:"deactivation_reason,omitempty"`
	DeactivationNotes     *string    `gorm:"type:text" json:"deactivation_notes,omitempty"`
	SuspensionUntil       *time.Time `json:"suspension_until,omitempty"`
	ReactivationRequested bool       `gorm:"default:false" json:"reactivation_requested,omitempty"`
	ContestationDeadline  *time.Time `json:"contestation_deadline,omitempty"`
}
