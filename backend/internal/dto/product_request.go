package dto

import (
	"github.com/lib/pq"
)

// ProductFormDTO represents the expected payload for creating a product.
type ProductFormDTO struct {
	Name          string         `form:"name" binding:"required"`
	Brand         string         `form:"brand" binding:"required"`
	Weight        float64        `form:"weight" binding:"required"`
	Description   string         `form:"description"`
	Price         float64        `form:"price" binding:"required"`
	Category      string         `form:"category" binding:"required"`
	StockQuantity int            `form:"stock_quantity" binding:"required,gte=0"`
	Accords       pq.StringArray `form:"accords"`
	Occasions     pq.StringArray `form:"occasions"`
	Seasons       pq.StringArray `form:"seasons"`
	Intensity     string         `form:"intensity"`
	Gender        string         `form:"gender"`
	PriceRange    string         `form:"price_range"`
	NotesTop      pq.StringArray `form:"notes_top"`
	NotesHeart    pq.StringArray `form:"notes_heart"`
	NotesBase     pq.StringArray `form:"notes_base"`
}

// UpdateProductRequest represents the payload for updating a product.
// @Description Product update request
type UpdateProductRequest struct {
	Name          *string         `json:"name,omitempty" example:"Sauvage Elixir"`
	Brand         *string         `json:"brand,omitempty" example:"Dior"`
	Weight        *float64        `json:"weight,omitempty" example:"60.0"`
	Description   *string         `json:"description,omitempty" example:"An intense and spicy fragrance"`
	Price         *float64        `json:"price,omitempty" example:"399.99"`
	Category      *string         `json:"category,omitempty" example:"Eau de Parfum"`
	StockQuantity *int            `json:"stock_quantity,omitempty" example:"25"`
	Accords       *pq.StringArray `json:"accords,omitempty"`
	Occasions     *pq.StringArray `json:"occasions,omitempty"`
	Seasons       *pq.StringArray `json:"seasons,omitempty"`
	Intensity     *string         `json:"intensity,omitempty"`
	Gender        *string         `json:"gender,omitempty"`
	PriceRange    *string         `json:"price_range,omitempty"`
	NotesTop      *pq.StringArray `json:"notes_top,omitempty"`
	NotesHeart    *pq.StringArray `json:"notes_heart,omitempty"`
	NotesBase     *pq.StringArray `json:"notes_base,omitempty"`
}
