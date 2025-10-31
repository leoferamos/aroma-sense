package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/leoferamos/aroma-sense/internal/dto"
	"github.com/leoferamos/aroma-sense/internal/model"
	"github.com/leoferamos/aroma-sense/internal/repository"
)

type OrderService interface {
	CreateOrderFromCart(userID string, req *dto.CreateOrderFromCartRequest) (*dto.OrderResponse, error)
}

type orderService struct {
	orderRepo   repository.OrderRepository
	cartRepo    repository.CartRepository
	productRepo repository.ProductRepository
}

func NewOrderService(orderRepo repository.OrderRepository, cartRepo repository.CartRepository, productRepo repository.ProductRepository) OrderService {
	return &orderService{orderRepo, cartRepo, productRepo}
}

func (s *orderService) CreateOrderFromCart(userID string, req *dto.CreateOrderFromCartRequest) (*dto.OrderResponse, error) {
	cart, err := s.cartRepo.FindByUserID(userID)
	if err != nil || cart == nil || len(cart.Items) == 0 {
		return nil, errors.New("cart is empty")
	}

	var orderItems []model.OrderItem
	total := 0.0
	for _, cartItem := range cart.Items {
		product, err := s.productRepo.FindByID(cartItem.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product not found: %d", cartItem.ProductID)
		}
		if product.StockQuantity < cartItem.Quantity {
			return nil, fmt.Errorf("insufficient stock for product: %s", product.Name)
		}
		itemSubtotal := float64(cartItem.Quantity) * product.Price
		orderItems = append(orderItems, model.OrderItem{
			ProductID:       product.ID,
			Quantity:        cartItem.Quantity,
			PriceAtPurchase: product.Price,
			Subtotal:        itemSubtotal,
		})
		total += itemSubtotal
	}

	order := &model.Order{
		UserID:          userID,
		TotalAmount:     total,
		Status:          model.OrderStatusPending,
		ShippingAddress: req.ShippingAddress,
		PaymentMethod:   model.PaymentMethod(req.PaymentMethod),
		Items:           orderItems,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.orderRepo.Create(order); err != nil {
		return nil, err
	}

	// Deduct stock
	for _, cartItem := range cart.Items {
		if err := s.productRepo.DecrementStock(cartItem.ProductID, cartItem.Quantity); err != nil {
			return nil, fmt.Errorf("failed to update stock for product %d", cartItem.ProductID)
		}
	}

	// Clear cart
	if err := s.cartRepo.ClearCartItems(cart.ID); err != nil {
		return nil, errors.New("order created, but failed to clear cart")
	}

	// Map order items to response
	items := make([]dto.OrderItemResponse, len(order.Items))
	for i, item := range order.Items {
		items[i] = dto.OrderItemResponse{
			ID:              item.ID,
			ProductID:       item.ProductID,
			Quantity:        item.Quantity,
			PriceAtPurchase: item.PriceAtPurchase,
			Subtotal:        item.Subtotal,
		}
	}
	return &dto.OrderResponse{
		ID:              order.ID,
		UserID:          order.UserID,
		TotalAmount:     order.TotalAmount,
		Status:          string(order.Status),
		ShippingAddress: order.ShippingAddress,
		PaymentMethod:   string(order.PaymentMethod),
		Items:           items,
		ItemCount:       len(order.Items),
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       order.UpdatedAt,
	}, nil
}
