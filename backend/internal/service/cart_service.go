package service

import (
	"github.com/leoferamos/aroma-sense/internal/model"
	"github.com/leoferamos/aroma-sense/internal/repository"
)

// CartService defines the interface for cart-related business logic
type CartService interface {
	CreateCartForUser(userID string) error
	GetCartByUserID(userID string) (*model.Cart, error)
}

type cartService struct {
	repo repository.CartRepository
}

// NewCartService creates a new instance of CartService
func NewCartService(repo repository.CartRepository) CartService {
	return &cartService{repo: repo}
}

// CreateCartForUser creates a new empty cart for a user
func (s *cartService) CreateCartForUser(userID string) error {
	// Check if cart already exists
	_, err := s.repo.FindByUserID(userID)
	if err == nil {
		return nil
	}

	cart := model.Cart{
		UserID: userID,
		Items:  []model.CartItem{},
	}
	return s.repo.Create(&cart)
}

// GetCartByUserID retrieves a cart by user ID
func (s *cartService) GetCartByUserID(userID string) (*model.Cart, error) {
	return s.repo.FindByUserID(userID)
}
