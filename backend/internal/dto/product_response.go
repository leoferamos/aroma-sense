package dto

import "time"

type ProductResponse struct {
	ID            uint      `json:"id"`
	Name          string    `json:"name"`
	Brand         string    `json:"brand"`
	Weight        float64   `json:"weight"`
	Description   string    `json:"description"`
	Price         float64   `json:"price"`
	ImageURL      string    `json:"image_url"`
	Category      string    `json:"category"`
	Notes         string    `json:"notes"`
	StockQuantity int       `json:"stock_quantity"`
	CreatedAt     time.Time `json:"created_at"`
}
