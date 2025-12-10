package model

import "time"

// ReviewReport represents a user-submitted report against a review.
type ReviewReport struct {
	ID             string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ReviewID       string    `gorm:"type:uuid;not null;index" json:"review_id"`
	Review         *Review   `gorm:"foreignKey:ReviewID" json:"review,omitempty"`
	ReportedBy     string    `gorm:"type:uuid;not null;index" json:"reported_by"`
	Reporter       *User     `gorm:"foreignKey:ReportedBy;references:PublicID" json:"reporter,omitempty"`
	ReasonCategory string    `gorm:"type:varchar(32);not null" json:"reason_category"`
	ReasonText     string    `gorm:"type:varchar(500)" json:"reason_text"`
	Status         string    `gorm:"type:varchar(16);not null;default:'pending';index" json:"status"`
	CreatedAt      time.Time `gorm:"autoCreateTime;index" json:"created_at"`
}
