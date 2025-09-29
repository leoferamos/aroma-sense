package repository

import (
	"github.com/leoferamos/aroma-sense/internal/model"
	"gorm.io/gorm"
)

// CartRepository defines the interface for cart data access operations
type CartRepository interface {
	Create(cart *model.Cart) error
	FindByUserID(userID string) (*model.Cart, error)
	Update(cart *model.Cart) error
	Delete(id uint) error
}

type cartRepository struct {
	db *gorm.DB
}

// NewCartRepository creates a new instance of CartRepository
func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db: db}
}

// Create inserts a new cart into the database
func (r *cartRepository) Create(cart *model.Cart) error {
	return r.db.Create(cart).Error
}

// FindByUserID retrieves a cart by user ID with items and products preloaded
func (r *cartRepository) FindByUserID(userID string) (*model.Cart, error) {
	var cart model.Cart
	err := r.db.Where("user_id = ?", userID).
		Preload("Items").
		Preload("Items.Product").
		First(&cart).Error
	return &cart, err
}

// Update modifies an existing cart
func (r *cartRepository) Update(cart *model.Cart) error {
	return r.db.Save(cart).Error
}

// Delete removes a cart by its ID
func (r *cartRepository) Delete(id uint) error {
	return r.db.Delete(&model.Cart{}, id).Error
}
