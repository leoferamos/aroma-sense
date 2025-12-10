package service

import (
	"context"
	"fmt"
	"time"

	"github.com/leoferamos/aroma-sense/internal/apperror"
	"github.com/leoferamos/aroma-sense/internal/dto"
	"github.com/leoferamos/aroma-sense/internal/model"
	"github.com/leoferamos/aroma-sense/internal/repository"
	"github.com/leoferamos/aroma-sense/internal/validation"
)

type OrderService interface {
	CreateOrderFromCart(userID string, req *dto.CreateOrderFromCartRequest) (*dto.OrderResponse, error)
	AdminListOrders(status *string, startDate *time.Time, endDate *time.Time, page int, perPage int) (*dto.AdminOrdersResponse, error)
	GetOrdersByUser(userID string) ([]dto.OrderResponse, error)
}

type orderService struct {
	orderRepo   repository.OrderRepository
	cartRepo    repository.CartRepository
	productRepo repository.ProductRepository
	shippingSvc ShippingService
}

func NewOrderService(orderRepo repository.OrderRepository, cartRepo repository.CartRepository, productRepo repository.ProductRepository, shippingSvc ShippingService) OrderService {
	return &orderService{orderRepo: orderRepo, cartRepo: cartRepo, productRepo: productRepo, shippingSvc: shippingSvc}
}

func (s *orderService) CreateOrderFromCart(userID string, req *dto.CreateOrderFromCartRequest) (*dto.OrderResponse, error) {
	cart, err := s.cartRepo.FindByUserID(userID)
	if err != nil || cart == nil || len(cart.Items) == 0 {
		return nil, apperror.NewCodeMessage("cart_empty", "cart is empty")
	}

	var orderItems []model.OrderItem
	total := 0.0
	for _, cartItem := range cart.Items {
		product, err := s.productRepo.FindByID(cartItem.ProductID)
		if err != nil {
			return nil, apperror.NewDomain(fmt.Errorf("product not found: %d", cartItem.ProductID), "product_not_found", "product not found")
		}
		if product.StockQuantity < cartItem.Quantity {
			return nil, apperror.NewDomain(fmt.Errorf("insufficient stock for product: %s", product.Name), "insufficient_stock", "insufficient stock")
		}
		itemSubtotal := float64(cartItem.Quantity) * product.Price
		orderItems = append(orderItems, model.OrderItem{
			ProductID:       product.ID,
			ProductSlug:     product.Slug,
			ProductName:     product.Name,
			ProductImageURL: product.ImageURL,
			Quantity:        cartItem.Quantity,
			PriceAtPurchase: product.Price,
			Subtotal:        itemSubtotal,
		})
		total += itemSubtotal
	}

	// Initialize order
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

	// Validate shipping selection against fresh quotes and persist shipping fields
	if req.ShippingSelection != nil {
		cep := validation.ExtractCEPFromString(req.ShippingAddress)
		if cep == "" {
			return nil, apperror.NewCodeMessage("invalid_postal_code", "invalid destination postal code")
		}

		if s.shippingSvc == nil {
			return nil, apperror.NewCodeMessage("provider_unavailable", "shipping provider not configured")
		}

		// Re-quote to validate selection.
		options, err := s.shippingSvc.CalculateOptions(context.Background(), userID, cep)
		if err != nil {
			return nil, err
		}

		var matched *dto.ShippingOption
		for i := range options {
			if options[i].Carrier == req.ShippingSelection.Carrier && options[i].ServiceCode == req.ShippingSelection.ServiceCode {
				matched = &options[i]
				break
			}
		}
		if matched == nil {
			return nil, apperror.NewCodeMessage("invalid_shipping_selection", "invalid shipping selection")
		}

		order.ShippingCarrier = matched.Carrier
		order.ShippingServiceCode = matched.ServiceCode
		order.ShippingPrice = matched.Price
		if matched.EstimatedDays > 0 {
			eta := time.Now().Add(time.Hour * 24 * time.Duration(matched.EstimatedDays))
			order.ShippingEstimatedDelivery = &eta
		}
		// Update order total to include shipping
		order.TotalAmount += matched.Price
	}

	if err := s.orderRepo.Create(order); err != nil {
		return nil, err
	}

	// Clear cart
	if err := s.cartRepo.ClearCartItems(cart.ID); err != nil {
		return nil, apperror.NewCodeMessage("cart_clear_failed", "order created, but failed to clear cart")
	}

	// Map order items to response
	items := make([]dto.OrderItemResponse, len(order.Items))
	for i, item := range order.Items {
		items[i] = dto.OrderItemResponse{
			ProductSlug:     item.ProductSlug,
			ProductName:     item.ProductName,
			ProductImageURL: item.ProductImageURL,
			Quantity:        item.Quantity,
			PriceAtPurchase: item.PriceAtPurchase,
			Subtotal:        item.Subtotal,
		}
	}
	return &dto.OrderResponse{
		PublicID:                  order.PublicID,
		TotalAmount:               order.TotalAmount,
		Status:                    string(order.Status),
		ShippingAddress:           order.ShippingAddress,
		PaymentMethod:             string(order.PaymentMethod),
		ShippingPrice:             order.ShippingPrice,
		ShippingCarrier:           order.ShippingCarrier,
		ShippingServiceCode:       order.ShippingServiceCode,
		ShippingEstimatedDelivery: order.ShippingEstimatedDelivery,
		Items:                     items,
		ItemCount:                 len(order.Items),
		CreatedAt:                 order.CreatedAt,
		UpdatedAt:                 order.UpdatedAt,
	}, nil
}

// AdminListOrders returns orders for admin listing with pagination and stats
func (s *orderService) AdminListOrders(status *string, startDate *time.Time, endDate *time.Time, page int, perPage int) (*dto.AdminOrdersResponse, error) {
	orders, totalCount, totalRevenue, err := s.orderRepo.ListOrders(status, startDate, endDate, page, perPage)
	if err != nil {
		return nil, err
	}

	// Map to DTOs
	var items []dto.AdminOrderItem
	for _, o := range orders {
		items = append(items, dto.AdminOrderItem{
			ID:          o.ID,
			PublicID:    o.PublicID,
			UserID:      o.UserID,
			TotalAmount: o.TotalAmount,
			Status:      string(o.Status),
			CreatedAt:   o.CreatedAt,
		})
	}

	totalPages := 0
	if perPage > 0 {
		totalPages = int((totalCount + int64(perPage) - 1) / int64(perPage))
	}

	resp := &dto.AdminOrdersResponse{}
	resp.Orders = items
	resp.Meta.Pagination = dto.PaginationMeta{
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
		TotalCount: int(totalCount),
	}
	avg := 0.0
	if totalCount > 0 {
		avg = totalRevenue / float64(totalCount)
	}
	resp.Meta.Stats = dto.StatsMeta{
		TotalRevenue:      totalRevenue,
		AverageOrderValue: avg,
	}

	return resp, nil
}

// GetOrdersByUser returns the orders for a given user
func (s *orderService) GetOrdersByUser(userID string) ([]dto.OrderResponse, error) {
	orders, err := s.orderRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	resp := make([]dto.OrderResponse, 0, len(orders))
	for _, o := range orders {
		items := make([]dto.OrderItemResponse, len(o.Items))
		for i, it := range o.Items {
			item := dto.OrderItemResponse{
				ProductSlug:     it.ProductSlug,
				ProductName:     it.ProductName,
				ProductImageURL: it.ProductImageURL,
				Quantity:        it.Quantity,
				PriceAtPurchase: it.PriceAtPurchase,
				Subtotal:        it.Subtotal,
			}
			items[i] = item
		}

		resp = append(resp, dto.OrderResponse{
			PublicID:                  o.PublicID,
			TotalAmount:               o.TotalAmount,
			Status:                    string(o.Status),
			ShippingAddress:           o.ShippingAddress,
			PaymentMethod:             string(o.PaymentMethod),
			ShippingPrice:             o.ShippingPrice,
			ShippingCarrier:           o.ShippingCarrier,
			ShippingServiceCode:       o.ShippingServiceCode,
			ShippingEstimatedDelivery: o.ShippingEstimatedDelivery,
			ShippingTracking:          o.ShippingTracking,
			ShippingStatus:            o.ShippingStatus,
			Items:                     items,
			ItemCount:                 len(items),
			CreatedAt:                 o.CreatedAt,
			UpdatedAt:                 o.UpdatedAt,
		})
	}

	if resp == nil {
		resp = []dto.OrderResponse{}
	}
	return resp, nil
}
