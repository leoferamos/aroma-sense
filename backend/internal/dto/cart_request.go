package dto

// CartResponse represents the cart data returned to the client
type CartResponse struct {
	Items     []CartItemResponse `json:"items"`
	Total     float64            `json:"total"`
	ItemCount int                `json:"item_count"`
}

// CartItemResponse represents a cart item returned to the client
type CartItemResponse struct {
	Product  *ProductResponse `json:"product,omitempty"`
	Quantity int              `json:"quantity"`
	Price    float64          `json:"price"`
	Total    float64          `json:"total"`
}

// AddToCartRequest represents the payload for adding an item to cart
type AddToCartRequest struct {
	ProductSlug string `json:"product_slug" binding:"required"`
	Quantity    int    `json:"quantity" binding:"required,min=1"`
}

// UpdateCartItemRequest represents the payload for updating cart item quantity
type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" binding:"required,min=0"`
}
