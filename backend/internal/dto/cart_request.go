package dto

import "time"

// CartResponse represents the cart data returned to the client
type CartResponse struct {
	ID        uint               `json:"id"`
	UserID    string             `json:"user_id"`
	Items     []CartItemResponse `json:"items"`
	Total     float64            `json:"total"`
	ItemCount int                `json:"item_count"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

// CartItemResponse represents a cart item returned to the client
type CartItemResponse struct {
	ID        uint             `json:"id"`
	CartID    uint             `json:"cart_id"`
	ProductID uint             `json:"product_id"`
	Product   *ProductResponse `json:"product,omitempty"`
	Quantity  int              `json:"quantity"`
	Price     float64          `json:"price"`
	Total     float64          `json:"total"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

// AddToCartRequest represents the payload for adding an item to cart
type AddToCartRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,min=1"`
}

// UpdateCartItemRequest represents the payload for updating cart item quantity
type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}
