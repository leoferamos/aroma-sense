package model

import (
	"time"

	"gorm.io/datatypes"
)

// PaymentStatus represents the status of a payment intent.
type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "pending"
	PaymentStatusProcessing PaymentStatus = "processing"
	PaymentStatusSucceeded  PaymentStatus = "succeeded"
	PaymentStatusFailed     PaymentStatus = "failed"
	PaymentStatusCanceled   PaymentStatus = "canceled"
)

// Payment stores gateway intent information for reconciliation.
type Payment struct {
	ID            uint              `gorm:"primaryKey" json:"-"`
	IntentID      string            `gorm:"size:255;not null;uniqueIndex" json:"intent_id"`
	Provider      string            `gorm:"size:50;not null" json:"provider"`
	UserID        string            `gorm:"size:255;not null;index" json:"user_id"`
	OrderPublicID *string           `gorm:"type:uuid;index" json:"order_public_id,omitempty"`
	AmountCents   int64             `gorm:"not null" json:"amount_cents"`
	Currency      string            `gorm:"size:10;not null" json:"currency"`
	Status        PaymentStatus     `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	Metadata      datatypes.JSONMap `gorm:"type:jsonb" json:"metadata,omitempty"`
	ErrorCode     string            `gorm:"size:100" json:"error_code,omitempty"`
	ErrorMessage  string            `gorm:"type:text" json:"error_message,omitempty"`
	CreatedAt     time.Time         `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time         `gorm:"autoUpdateTime" json:"updated_at"`
}
