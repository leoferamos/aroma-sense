package model

import "time"

// Cart represents a user's shopping cart.
type Cart struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	UserID    string     `gorm:"size:255;not null;uniqueIndex" json:"user_id"`
	Items     []CartItem `gorm:"foreignKey:CartID" json:"items"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

// CartItem represents an item in a shopping cart.
type CartItem struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CartID    uint      `gorm:"not null" json:"cart_id"`
	ProductID uint      `gorm:"not null" json:"product_id"`
	Product   *Product  `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Quantity  int       `gorm:"not null;default:1" json:"quantity"`
	Price     float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
