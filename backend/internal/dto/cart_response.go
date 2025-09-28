package dto

import "time"

// CartSummaryResponse represents a summary of the cart
type CartSummaryResponse struct {
	Cart    CartResponse `json:"cart"`
	Summary CartSummary  `json:"summary"`
}

// CartSummary provides additional cart information
type CartSummary struct {
	ItemCount   int       `json:"item_count"`
	TotalAmount float64   `json:"total_amount"`
	UniqueItems int       `json:"unique_items"`
	LastUpdated time.Time `json:"last_updated"`
}
