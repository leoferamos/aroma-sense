package model

import "time"

// Product represents a fragance/product in the catalog.
type Product struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Name          string    `gorm:"size:128;not null" json:"name"`
	Brand         string    `gorm:"size:64;not null" json:"brand"`
	Weight        float64   `gorm:"not null" json:"weight"`
	Description   string    `gorm:"type:text" json:"description"`
	Price         float64   `gorm:"not null" json:"price"`
	ImageURL      string    `gorm:"size:256" json:"image_url"`
	ThumbnailURL  string    `gorm:"size:256" json:"thumbnail_url"`
	Category      string    `gorm:"size:64;not null" json:"category"`
	Notes         string    `gorm:"type:text" json:"notes"`
	StockQuantity int       `gorm:"not null" json:"stock_quantity"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
