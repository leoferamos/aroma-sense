package service

import (
	"context"
	"errors"

	"github.com/leoferamos/aroma-sense/internal/dto"
	"github.com/leoferamos/aroma-sense/internal/model"
	"github.com/leoferamos/aroma-sense/internal/repository"
)

// CartService defines the interface for cart-related business logic
type CartService interface {
	CreateCartForUser(userID string) error
	GetCartByUserID(userID string) (*model.Cart, error)
	GetCartResponse(userID string) (*dto.CartResponse, error)
	AddItemToCart(userID string, productID uint, quantity int) (*dto.CartResponse, error)
}

type cartService struct {
	repo           repository.CartRepository
	productService ProductService
}

// NewCartService creates a new instance of CartService
func NewCartService(repo repository.CartRepository, productService ProductService) CartService {
	return &cartService{repo: repo, productService: productService}
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

// GetCartResponse retrieves a cart response by user ID
func (s *cartService) GetCartResponse(userID string) (*dto.CartResponse, error) {
	cart, err := s.GetCartByUserID(userID)
	if err != nil {
		return nil, err
	}

	cartResponse := &dto.CartResponse{
		ID:        cart.ID,
		UserID:    cart.UserID,
		Items:     []dto.CartItemResponse{},
		Total:     0.0,
		ItemCount: 0,
		CreatedAt: cart.CreatedAt,
		UpdatedAt: cart.UpdatedAt,
	}

	// Convert cart items and calculate totals
	for _, item := range cart.Items {
		itemTotal := item.Price * float64(item.Quantity)

		cartItemResponse := dto.CartItemResponse{
			ID:        item.ID,
			CartID:    item.CartID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
			Total:     itemTotal,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		}

		if item.Product != nil {
			cartItemResponse.Product = &dto.ProductResponse{
				ID:            item.Product.ID,
				Name:          item.Product.Name,
				Brand:         item.Product.Brand,
				Weight:        item.Product.Weight,
				Description:   item.Product.Description,
				Price:         item.Product.Price,
				ImageURL:      item.Product.ImageURL,
				Category:      item.Product.Category,
				Notes:         item.Product.Notes,
				StockQuantity: item.Product.StockQuantity,
				CreatedAt:     item.Product.CreatedAt,
				UpdatedAt:     item.Product.UpdatedAt,
			}
		}

		cartResponse.Items = append(cartResponse.Items, cartItemResponse)
		cartResponse.Total += itemTotal
		cartResponse.ItemCount += item.Quantity
	}

	return cartResponse, nil
}

// AddItemToCart adds an item to the user's cart or increases quantity if exists
func (s *cartService) AddItemToCart(userID string, productID uint, quantity int) (*dto.CartResponse, error) {
	// Validate product exists and get product data
	product, err := s.productService.GetProductByID(context.Background(), productID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	// Check stock availability
	if product.StockQuantity < quantity {
		return nil, errors.New("insufficient stock")
	}

	// Get user's cart
	cart, err := s.GetCartByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Check if item already exists in cart
	var existingItem *model.CartItem
	for i := range cart.Items {
		if cart.Items[i].ProductID == productID {
			existingItem = &cart.Items[i]
			break
		}
	}

	if existingItem != nil {
		// Check total stock for new quantity
		totalQuantity := existingItem.Quantity + quantity
		if product.StockQuantity < totalQuantity {
			return nil, errors.New("insufficient stock for requested quantity")
		}

		// Update quantity and save to database
		existingItem.Quantity = totalQuantity
		if err := s.repo.UpdateCartItem(existingItem); err != nil {
			return nil, errors.New("failed to update cart item")
		}
	} else {
		// Create new cart item
		newItem := model.CartItem{
			CartID:    cart.ID,
			ProductID: productID,
			Quantity:  quantity,
			Price:     product.Price,
		}

		// Save to database
		if err := s.repo.CreateCartItem(&newItem); err != nil {
			return nil, errors.New("failed to add item to cart")
		}
	}

	// Return updated cart
	return s.GetCartResponse(userID)
}
