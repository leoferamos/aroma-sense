package model

import (
	"time"
)

// UserContestation represents a request to contest account deactivation
type UserContestation struct {
	ID          uint       `gorm:"primaryKey;column:id"`
	UserID      uint       `gorm:"not null;index;column:user_id"`
	Reason      string     `gorm:"type:text;not null;column:reason"`
	Status      string     `gorm:"type:varchar(16);not null;default:pending;index;column:status"`
	RequestedAt time.Time  `gorm:"not null;autoCreateTime;column:requested_at"`
	ReviewedAt  *time.Time `gorm:"column:reviewed_at"`
	ReviewedBy  *uint      `gorm:"column:reviewed_by"`
	ReviewNotes *string    `gorm:"type:text;column:review_notes"`
}
