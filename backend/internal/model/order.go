package model

import "time"

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
)

// PaymentMethod represents the payment method used for an order
type PaymentMethod string

const (
	PaymentMethodCreditCard PaymentMethod = "credit_card"
	PaymentMethodDebitCard  PaymentMethod = "debit_card"
	PaymentMethodPix        PaymentMethod = "pix"
	PaymentMethodBoleto     PaymentMethod = "boleto"
)

// Order represents a customer order.
type Order struct {
	ID              uint          `gorm:"primaryKey" json:"id"`
	UserID          string        `gorm:"size:255;not null;index" json:"user_id"`
	User            *User         `gorm:"foreignKey:UserID;references:PublicID" json:"user,omitempty"`
	TotalAmount     float64       `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	Status          OrderStatus   `gorm:"type:varchar(20);not null;default:'pending';index" json:"status"`
	ShippingAddress string        `gorm:"type:text;not null" json:"shipping_address"`
	PaymentMethod   PaymentMethod `gorm:"type:varchar(20);not null" json:"payment_method"`
	Items           []OrderItem   `gorm:"foreignKey:OrderID" json:"items"`
	CreatedAt       time.Time     `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt       time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
}

// OrderItem represents an item in an order.
type OrderItem struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	OrderID         uint      `gorm:"not null;index" json:"order_id"`
	ProductID       uint      `gorm:"not null;index" json:"product_id"`
	Product         *Product  `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Quantity        int       `gorm:"not null" json:"quantity"`
	PriceAtPurchase float64   `gorm:"type:decimal(10,2);not null" json:"price_at_purchase"`
	Subtotal        float64   `gorm:"type:decimal(10,2);not null" json:"subtotal"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

