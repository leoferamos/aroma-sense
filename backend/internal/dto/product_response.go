package dto

import "time"

// ProductResponse represents the product data returned to the client
// @Description Product information returned by the API
type ProductResponse struct {
	Name               string    `json:"name" example:"Sauvage"`
	Brand              string    `json:"brand" example:"Dior"`
	Weight             float64   `json:"weight" example:"100.0"`
	Description        string    `json:"description" example:"A fresh and woody fragrance"`
	Price              float64   `json:"price" example:"299.99"`
	ImageURL           string    `json:"image_url" example:"https://example.com/image.jpg"`
	ThumbnailURL       string    `json:"thumbnail_url,omitempty" example:"https://example.com/image_thumb.jpg"`
	Slug               string    `json:"slug,omitempty" example:"dior-sauvage"`
	Accords            []string  `json:"accords,omitempty" example:"[\"woody\",\"citrus\"]"`
	Occasions          []string  `json:"occasions,omitempty" example:"[\"work\",\"night out\"]"`
	Seasons            []string  `json:"seasons,omitempty" example:"[\"summer\",\"spring\"]"`
	Intensity          string    `json:"intensity,omitempty" example:"moderate"`
	Gender             string    `json:"gender,omitempty" example:"unisex"`
	PriceRange         string    `json:"price_range,omitempty" example:"premium"`
	NotesTop           []string  `json:"notes_top,omitempty" example:"[\"bergamot\"]"`
	NotesHeart         []string  `json:"notes_heart,omitempty" example:"[\"lavender\"]"`
	NotesBase          []string  `json:"notes_base,omitempty" example:"[\"ambroxan\"]"`
	Category           string    `json:"category" example:"Eau de Parfum"`
	StockQuantity      int       `json:"stock_quantity" example:"50"`
	CreatedAt          time.Time `json:"created_at" example:"2025-09-28T10:00:00Z"`
	UpdatedAt          time.Time `json:"updated_at" example:"2025-09-28T10:00:00Z"`
	CanReview          *bool     `json:"can_review,omitempty"`
	CannotReviewReason *string   `json:"cannot_review_reason,omitempty"`
}
