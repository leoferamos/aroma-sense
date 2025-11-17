package model

import "time"

// ReviewStatus represents moderation status for a review
type ReviewStatus string

const (
	ReviewStatusPublished ReviewStatus = "published"
	ReviewStatusHidden    ReviewStatus = "hidden"
	ReviewStatusFlagged   ReviewStatus = "flagged"
)

// Review represents a product review authored by a user
type Review struct {
	ID        string       `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ProductID uint         `gorm:"not null;index" json:"product_id"`
	Product   *Product     `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	UserID    string       `gorm:"type:uuid;not null;index" json:"user_id"`
	User      *User        `gorm:"foreignKey:UserID;references:PublicID" json:"user,omitempty"`
	Rating    int          `gorm:"not null;check:rating >= 1 AND rating <= 5" json:"rating"`
	Comment   string       `gorm:"type:text" json:"comment"`
	Status    ReviewStatus `gorm:"type:varchar(16);not null;default:'published';index" json:"status"`
	CreatedAt time.Time    `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt time.Time    `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt *time.Time   `gorm:"index" json:"-"`
}
