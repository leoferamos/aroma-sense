package model

import "time"

// ResetToken represents a token for password reset requests.
type ResetToken struct {
	ID        uint      `gorm:"primaryKey" json:"-"`
	Email     string    `gorm:"size:100;not null;index:idx_password_reset_tokens_email" json:"email"`
	Code      string    `gorm:"size:6;not null" json:"-"`
	ExpiresAt time.Time `gorm:"not null;index:idx_password_reset_tokens_expires_at" json:"expires_at"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// TableName specifies the table name for GORM.
func (ResetToken) TableName() string {
	return "password_reset_tokens"
}

// IsExpired checks if the token has expired.
func (t *ResetToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}
