package dto

import "time"

// OrderResponse represents the order data returned to the client
type OrderResponse struct {
	ID              uint                `json:"id"`
	UserID          string              `json:"user_id"`
	TotalAmount     float64             `json:"total_amount"`
	Status          string              `json:"status"`
	ShippingAddress string              `json:"shipping_address"`
	PaymentMethod   string              `json:"payment_method"`
	Items           []OrderItemResponse `json:"items"`
	ItemCount       int                 `json:"item_count"`
	CreatedAt       time.Time           `json:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at"`
}

// OrderItemResponse represents an order item returned to the client
type OrderItemResponse struct {
	ID              uint             `json:"id"`
	ProductID       uint             `json:"product_id"`
	Product         *ProductResponse `json:"product,omitempty"`
	Quantity        int              `json:"quantity"`
	PriceAtPurchase float64          `json:"price_at_purchase"`
	Subtotal        float64          `json:"subtotal"`
}

