package service

import (
	"context"
	"fmt"

	"github.com/leoferamos/aroma-sense/internal/apperror"
	"github.com/leoferamos/aroma-sense/internal/dto"
	"github.com/leoferamos/aroma-sense/internal/model"
	"github.com/leoferamos/aroma-sense/internal/repository"
)

// CartService defines the interface for cart-related business logic
type CartService interface {
	CreateCartForUser(userID string) error
	GetCartByUserID(userID string) (*model.Cart, error)
	GetCartResponse(userID string) (*dto.CartResponse, error)
	AddItemToCart(userID string, productSlug string, quantity int) (*dto.CartResponse, error)
	UpdateItemQuantity(userID string, itemID uint, quantity int) (*dto.CartResponse, error)
	UpdateItemQuantityBySlug(userID string, productSlug string, quantity int) (*dto.CartResponse, error)
	RemoveItem(userID string, itemID uint) (*dto.CartResponse, error)
	RemoveItemBySlug(userID string, productSlug string) (*dto.CartResponse, error)
	ClearCart(userID string) (*dto.CartResponse, error)
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
		Items:     []dto.CartItemResponse{},
		Total:     0.0,
		ItemCount: 0,
	}

	// Convert cart items and calculate totals
	for _, item := range cart.Items {
		itemTotal := item.Price * float64(item.Quantity)

		cartItemResponse := dto.CartItemResponse{
			Quantity: item.Quantity,
			Price:    item.Price,
			Total:    itemTotal,
		}

		if item.Product != nil {
			cartItemResponse.Product = &dto.ProductResponse{
				Name:          item.Product.Name,
				Brand:         item.Product.Brand,
				Weight:        item.Product.Weight,
				Description:   item.Product.Description,
				Price:         item.Product.Price,
				ImageURL:      item.Product.ImageURL,
				ThumbnailURL:  item.Product.ThumbnailURL,
				Slug:          item.Product.Slug,
				Category:      item.Product.Category,
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
func (s *cartService) AddItemToCart(userID string, productSlug string, quantity int) (*dto.CartResponse, error) {
	// Get product ID by slug
	productID, err := s.productService.GetProductIDBySlug(context.Background(), productSlug)
	if err != nil {
		return nil, apperror.NewCodeMessage("product_not_found", "product not found")
	}

	// Validate product exists and get product data
	product, err := s.productService.GetProductByID(context.Background(), productID)
	if err != nil {
		return nil, apperror.NewCodeMessage("product_not_found", "product not found")
	}

	// Check stock availability
	if product.StockQuantity <= 0 {
		return nil, apperror.NewCodeMessage("insufficient_stock", "product out of stock")
	}
	if product.StockQuantity < quantity {
		return nil, apperror.NewCodeMessage("insufficient_stock", fmt.Sprintf("insufficient stock - only %d items available", product.StockQuantity))
	}

	// Get user's cart
	cart, err := s.GetCartByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Check if item already exists in cart
	var existingItemID uint
	itemExists := false
	for _, item := range cart.Items {
		if item.ProductID == productID {
			existingItemID = item.ID
			itemExists = true
			// Calculate new total quantity and validate stock
			newQuantity := item.Quantity + quantity
			if product.StockQuantity <= 0 {
				return nil, apperror.NewCodeMessage("insufficient_stock", "product out of stock")
			}
			if product.StockQuantity < newQuantity {
				return nil, apperror.NewCodeMessage("insufficient_stock", fmt.Sprintf("insufficient stock - only %d items available, you already have %d in cart", product.StockQuantity, item.Quantity))
			}
			break
		}
	}

	if itemExists {
		// Get current quantity and add the new quantity
		currentItem, err := s.repo.FindCartItemByID(existingItemID)
		if err != nil {
			return nil, apperror.NewDomain(err, "cart_item_not_found", "cart item not found")
		}
		newQuantity := currentItem.Quantity + quantity
		return s.UpdateItemQuantity(userID, existingItemID, newQuantity)
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
			return nil, apperror.NewDomain(err, "cart_update_failed", "failed to add item to cart")
		}
	}

	// Return updated cart
	return s.GetCartResponse(userID)
}

// UpdateItemQuantity updates the quantity of a specific cart item
func (s *cartService) UpdateItemQuantity(userID string, itemID uint, quantity int) (*dto.CartResponse, error) {
	// Get the cart item
	cartItem, err := s.repo.FindCartItemByID(itemID)
	if err != nil {
		return nil, apperror.NewCodeMessage("cart_item_not_found", "cart item not found")
	}

	// Verify the item belongs to the user's cart
	cart, err := s.GetCartByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Check if the item belongs to this user's cart
	itemBelongsToUser := false
	for _, item := range cart.Items {
		if item.ID == itemID {
			itemBelongsToUser = true
			break
		}
	}

	if !itemBelongsToUser {
		return nil, apperror.NewCodeMessage("cart_item_not_found", "cart item not found in user's cart")
	}

	// If quantity is 0, remove the item
	if quantity == 0 {
		if err := s.repo.DeleteCartItem(itemID); err != nil {
			return nil, apperror.NewDomain(err, "cart_update_failed", "failed to remove cart item")
		}
	} else {
		// Validate stock availability for the new quantity
		product, err := s.productService.GetProductByID(context.Background(), cartItem.ProductID)
		if err != nil {
			return nil, apperror.NewCodeMessage("product_not_found", "product not found")
		}

		if product.StockQuantity <= 0 {
			return nil, apperror.NewCodeMessage("insufficient_stock", "product out of stock")
		}
		if product.StockQuantity < quantity {
			return nil, apperror.NewCodeMessage("insufficient_stock", fmt.Sprintf("insufficient stock - only %d items available", product.StockQuantity))
		}

		// Update the quantity
		cartItem.Quantity = quantity
		if err := s.repo.UpdateCartItem(cartItem); err != nil {
			return nil, apperror.NewDomain(err, "cart_update_failed", "failed to update cart item quantity")
		}
	}

	// Return updated cart
	return s.GetCartResponse(userID)
}

// RemoveItem removes an item from the user's cart
func (s *cartService) RemoveItem(userID string, itemID uint) (*dto.CartResponse, error) {
	return s.UpdateItemQuantity(userID, itemID, 0)
}

// ClearCart removes all items from the user's cart
func (s *cartService) ClearCart(userID string) (*dto.CartResponse, error) {
	// Get user's cart
	cart, err := s.GetCartByUserID(userID)
	if err != nil {
		return nil, err
	}

	if err := s.repo.ClearCartItems(cart.ID); err != nil {
		return nil, apperror.NewCodeMessage("cart_clear_failed", "failed to clear cart")
	}

	// Return empty cart response
	return s.GetCartResponse(userID)
}

// UpdateItemQuantityBySlug updates the quantity of a cart item by product slug
func (s *cartService) UpdateItemQuantityBySlug(userID string, productSlug string, quantity int) (*dto.CartResponse, error) {
	// Get product ID by slug
	productID, err := s.productService.GetProductIDBySlug(context.Background(), productSlug)
	if err != nil {
		return nil, apperror.NewCodeMessage("product_not_found", "product not found")
	}

	// Get user's cart
	cart, err := s.GetCartByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Find the cart item by product ID
	var itemID uint
	found := false
	for _, item := range cart.Items {
		if item.ProductID == productID {
			itemID = item.ID
			found = true
			break
		}
	}

	if !found {
		return nil, apperror.NewCodeMessage("cart_item_not_found", "cart item not found in user's cart")
	}

	// Use existing UpdateItemQuantity method
	return s.UpdateItemQuantity(userID, itemID, quantity)
}

// RemoveItemBySlug removes a cart item by product slug
func (s *cartService) RemoveItemBySlug(userID string, productSlug string) (*dto.CartResponse, error) {
	// Get product ID by slug
	productID, err := s.productService.GetProductIDBySlug(context.Background(), productSlug)
	if err != nil {
		return nil, apperror.NewCodeMessage("product_not_found", "product not found")
	}

	// Get user's cart
	cart, err := s.GetCartByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Find the cart item by product ID
	var itemID uint
	found := false
	for _, item := range cart.Items {
		if item.ProductID == productID {
			itemID = item.ID
			found = true
			break
		}
	}

	if !found {
		return nil, apperror.NewCodeMessage("cart_item_not_found", "cart item not found in user's cart")
	}

	// Use existing RemoveItem method
	return s.RemoveItem(userID, itemID)
}
