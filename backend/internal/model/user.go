package model

import "time"

// User represents the database entity for a registered user.
type User struct {
	ID           uint      `gorm:"primaryKey" json:"-"`
	PublicID     string    `gorm:"type:uuid;not null;uniqueIndex" json:"public_id"`
	Username     string    `gorm:"size:64;not null;unique" json:"username"`
	Email        string    `gorm:"size:128;not null;unique" json:"email"`
	PasswordHash string    `gorm:"size:256;not null" json:"-"`
	Role         string    `gorm:"size:16;not null;default:client" json:"role"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
}
