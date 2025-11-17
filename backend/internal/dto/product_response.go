package dto

import "time"

// ProductResponse represents the product data returned to the client
// @Description Product information returned by the API
type ProductResponse struct {
	ID                 uint      `json:"id" example:"1"`
	Name               string    `json:"name" example:"Sauvage"`
	Brand              string    `json:"brand" example:"Dior"`
	Weight             float64   `json:"weight" example:"100.0"`
	Description        string    `json:"description" example:"A fresh and woody fragrance"`
	Price              float64   `json:"price" example:"299.99"`
	ImageURL           string    `json:"image_url" example:"https://example.com/image.jpg"`
	ThumbnailURL       string    `json:"thumbnail_url,omitempty" example:"https://example.com/image_thumb.jpg"`
	Category           string    `json:"category" example:"Eau de Parfum"`
	Notes              string    `json:"notes" example:"bergamot, pepper, ambroxan"`
	StockQuantity      int       `json:"stock_quantity" example:"50"`
	CreatedAt          time.Time `json:"created_at" example:"2025-09-28T10:00:00Z"`
	UpdatedAt          time.Time `json:"updated_at" example:"2025-09-28T10:00:00Z"`
	CanReview          *bool     `json:"can_review,omitempty"`
	CannotReviewReason *string   `json:"cannot_review_reason,omitempty"`
}
