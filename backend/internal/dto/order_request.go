package dto

// CreateOrderFromCartRequest represents the payload for creating an order from the entire cart
type CreateOrderFromCartRequest struct {
	ShippingAddress string `json:"shipping_address" binding:"required" example:"Rua Example, 123, São Paulo - SP, 01234-567"`
	PaymentMethod   string `json:"payment_method" binding:"required,oneof=credit_card debit_card pix boleto" example:"pix"`
}

// CreateOrderDirectRequest represents the payload for buying a single product directly
type CreateOrderDirectRequest struct {
	ProductID       uint   `json:"product_id" binding:"required" example:"5"`
	Quantity        int    `json:"quantity" binding:"required,min=1" example:"2"`
	ShippingAddress string `json:"shipping_address" binding:"required" example:"Rua Example, 123, São Paulo - SP, 01234-567"`
	PaymentMethod   string `json:"payment_method" binding:"required,oneof=credit_card debit_card pix boleto" example:"pix"`
}

// UpdateOrderStatusRequest represents the payload for updating order status 
type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending processing shipped delivered cancelled" example:"processing"`
}



