package model

import (
	"time"

	"github.com/lib/pq"
)

// Product represents a product in the catalog.
type Product struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"size:128;not null" json:"name"`
	Brand        string         `gorm:"size:64;not null" json:"brand"`
	Weight       float64        `gorm:"not null" json:"weight"`
	Description  string         `gorm:"type:text" json:"description"`
	Price        float64        `gorm:"not null" json:"price"`
	ImageURL     string         `gorm:"size:256" json:"image_url"`
	ThumbnailURL string         `gorm:"size:256" json:"thumbnail_url"`
	Slug         string         `gorm:"size:128" json:"slug,omitempty"`
	Accords      pq.StringArray `gorm:"type:text[]" json:"accords,omitempty"`
	Occasions    pq.StringArray `gorm:"type:text[]" json:"occasions,omitempty"`
	Seasons      pq.StringArray `gorm:"type:text[]" json:"seasons,omitempty"`
	Intensity    string         `gorm:"size:16" json:"intensity,omitempty"`
	Gender       string         `gorm:"size:16" json:"gender,omitempty"`
	PriceRange   string         `gorm:"size:16" json:"price_range,omitempty"`
	NotesTop     pq.StringArray `gorm:"type:text[]" json:"notes_top,omitempty"`
	NotesHeart   pq.StringArray `gorm:"type:text[]" json:"notes_heart,omitempty"`
	NotesBase    pq.StringArray `gorm:"type:text[]" json:"notes_base,omitempty"`

	Category      string    `gorm:"size:64;not null" json:"category"`
	StockQuantity int       `gorm:"not null" json:"stock_quantity"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
