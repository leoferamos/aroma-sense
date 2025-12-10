package service

import (
	"context"
	"fmt"

	"github.com/leoferamos/aroma-sense/internal/apperror"
	"github.com/leoferamos/aroma-sense/internal/dto"
	"github.com/leoferamos/aroma-sense/internal/model"
	"github.com/leoferamos/aroma-sense/internal/repository"
	"gorm.io/datatypes"
)

// PaymentProvider defines the operations required from a payment gateway client.
type PaymentProvider interface {
	CreatePaymentIntent(ctx context.Context, params PaymentIntentParams) (*PaymentIntentResult, error)
	ParseWebhook(payload []byte, signature string) (*PaymentWebhookPayload, error)
}

// PaymentIntentParams are the normalized params sent to the provider.
type PaymentIntentParams struct {
	Amount        int64
	Currency      string
	CustomerEmail string
	Metadata      map[string]string
}

// PaymentIntentResult wraps the minimal fields we expose to callers.
type PaymentIntentResult struct {
	ID           string
	ClientSecret string
}

// PaymentWebhookPayload is a normalized view of provider webhook events.
type PaymentWebhookPayload struct {
	IntentID      string
	Status        string
	Amount        int64
	Currency      string
	CustomerEmail string
	Metadata      map[string]string
}

type PaymentService interface {
	CreateIntent(ctx context.Context, userID string, req *dto.CreatePaymentIntentRequest) (*PaymentIntentResult, error)
	HandleWebhook(ctx context.Context, payload []byte, signature string) (*PaymentWebhookPayload, error)
}

type paymentService struct {
	cartRepo    repository.CartRepository
	productRepo repository.ProductRepository
	orderRepo   repository.OrderRepository
	paymentRepo repository.PaymentRepository
	shippingSvc ShippingService
	provider    PaymentProvider
}

func NewPaymentService(cartRepo repository.CartRepository, productRepo repository.ProductRepository, orderRepo repository.OrderRepository, paymentRepo repository.PaymentRepository, shippingSvc ShippingService, provider PaymentProvider) PaymentService {
	return &paymentService{cartRepo: cartRepo, productRepo: productRepo, orderRepo: orderRepo, paymentRepo: paymentRepo, shippingSvc: shippingSvc, provider: provider}
}

// CreateIntent calculates totals from the user's cart and shipping selection, then delegates to the provider.
func (s *paymentService) CreateIntent(ctx context.Context, userID string, req *dto.CreatePaymentIntentRequest) (*PaymentIntentResult, error) {
	if req == nil {
		return nil, apperror.NewCodeMessage("invalid_request", "missing payload")
	}

	var amount int64
	metadata := map[string]string{
		"user_id": userID,
	}

	// If an order already exists, use its totals and metadata.
	if req.OrderPublicID != "" {
		order, err := s.orderRepo.FindByPublicIDWithItems(req.OrderPublicID)
		if err != nil {
			return nil, err
		}
		if order == nil || order.UserID != userID {
			return nil, apperror.NewCodeMessage("invalid_request", "order not found")
		}
		amount = int64(order.TotalAmount * 100)
		metadata["order_public_id"] = req.OrderPublicID
		if order.ShippingAddress != "" {
			metadata["shipping_address"] = order.ShippingAddress
		}
		if order.ShippingCarrier != "" {
			metadata["shipping_carrier"] = order.ShippingCarrier
		}
		if order.ShippingServiceCode != "" {
			metadata["shipping_service_code"] = order.ShippingServiceCode
		}
	} else {
		cart, err := s.cartRepo.FindByUserID(userID)
		if err != nil || cart == nil || len(cart.Items) == 0 {
			return nil, apperror.NewCodeMessage("cart_empty", "cart is empty")
		}

		// Compute subtotal and validate stock existence.
		total := 0.0
		for _, item := range cart.Items {
			product, err := s.productRepo.FindByID(item.ProductID)
			if err != nil {
				return nil, apperror.NewDomain(fmt.Errorf("product not found: %d", item.ProductID), "product_not_found", "product not found")
			}
			if product.StockQuantity < item.Quantity {
				return nil, apperror.NewDomain(fmt.Errorf("insufficient stock for product: %s", product.Name), "insufficient_stock", "insufficient stock")
			}
			total += float64(item.Quantity) * product.Price
		}

		if req.ShippingSelection != nil {
			if s.shippingSvc == nil {
				return nil, apperror.NewCodeMessage("provider_unavailable", "shipping provider not configured")
			}
			options, err := s.shippingSvc.CalculateOptions(ctx, userID, req.ShippingPostalCode())
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
			total += matched.Price
			metadata["shipping_carrier"] = req.ShippingSelection.Carrier
			metadata["shipping_service_code"] = req.ShippingSelection.ServiceCode
		}

		amount = int64(total * 100) // convert to cents
		if req.ShippingAddress != "" {
			metadata["shipping_address"] = req.ShippingAddress
		}
		metadata["order_hint"] = "cart"
	}

	if amount <= 0 {
		return nil, apperror.NewCodeMessage("invalid_amount", "amount must be positive")
	}

	params := PaymentIntentParams{
		Amount:        amount,
		Currency:      "brl",
		CustomerEmail: req.CustomerEmail,
		Metadata:      metadata,
	}

	result, err := s.provider.CreatePaymentIntent(ctx, params)
	if err != nil {
		return nil, err
	}

	if s.paymentRepo != nil {
		payment := &model.Payment{
			IntentID:    result.ID,
			Provider:    "stripe",
			UserID:      userID,
			AmountCents: amount,
			Currency:    params.Currency,
			Status:      model.PaymentStatusPending,
			Metadata:    toJSONMap(metadata),
		}
		if v, ok := metadata["order_public_id"]; ok && v != "" {
			payment.OrderPublicID = &v
		}
		if err := s.paymentRepo.Create(payment); err != nil {
			return nil, err
		}
	}

	return result, nil
}

// HandleWebhook validates provider signature and returns the normalized payload.
func (s *paymentService) HandleWebhook(ctx context.Context, payload []byte, signature string) (*PaymentWebhookPayload, error) {
	if s.provider == nil {
		return nil, apperror.NewCodeMessage("provider_unavailable", "payment provider not configured")
	}

	normalized, err := s.provider.ParseWebhook(payload, signature)
	if err != nil {
		return nil, apperror.NewDomain(err, "invalid_webhook", "invalid webhook payload")
	}

	status := toPaymentStatus(normalized.Status)
	metadata := normalized.Metadata

	if s.paymentRepo == nil {
		return normalized, nil
	}

	payment, err := s.paymentRepo.FindByIntentID(normalized.IntentID)
	if err != nil {
		return nil, err
	}

	orderPublicID := metadata["order_public_id"]

	if payment == nil {
		userID := metadata["user_id"]
		p := &model.Payment{
			IntentID:    normalized.IntentID,
			Provider:    "stripe",
			UserID:      userID,
			AmountCents: normalized.Amount,
			Currency:    normalized.Currency,
			Status:      status,
			Metadata:    toJSONMap(metadata),
		}
		if orderPublicID != "" {
			p.OrderPublicID = &orderPublicID
		}
		if err := s.paymentRepo.Create(p); err != nil {
			return nil, err
		}
		payment = p
	} else {
		if orderPublicID != "" && payment.OrderPublicID == nil {
			if err := s.paymentRepo.AttachOrderPublicID(payment.IntentID, orderPublicID); err == nil {
				payment.OrderPublicID = &orderPublicID
			}
		}
		// Idempotent update.
		if payment.Status != status {
			errCode := ""
			errMsg := ""
			if status == model.PaymentStatusFailed || status == model.PaymentStatusCanceled {
				errCode = "payment_failed"
				errMsg = "provider reported failure"
			}
			if err := s.paymentRepo.UpdateStatusByIntentID(payment.IntentID, status, errCode, errMsg); err != nil {
				return nil, err
			}
			payment.Status = status
		}
	}

	if s.orderRepo != nil {
		target := orderPublicID
		if target == "" && payment != nil && payment.OrderPublicID != nil {
			target = *payment.OrderPublicID
		}
		if target != "" {
			order, err := s.orderRepo.FindByPublicIDWithItems(target)
			if err != nil {
				return nil, err
			}
			if order == nil {
				return normalized, nil
			}
			switch status {
			case model.PaymentStatusSucceeded:
				if order.Status == model.OrderStatusPending {
					for _, item := range order.Items {
						if err := s.productRepo.DecrementStock(item.ProductID, item.Quantity); err != nil {
							return nil, err
						}
					}
				}
				if err := s.orderRepo.UpdateStatusByPublicID(target, model.OrderStatusProcessing); err != nil {
					return nil, err
				}
			case model.PaymentStatusFailed, model.PaymentStatusCanceled:
				if err := s.orderRepo.UpdateStatusByPublicID(target, model.OrderStatusCancelled); err != nil {
					return nil, err
				}
			}
		}
	}

	return normalized, nil
}

func toPaymentStatus(raw string) model.PaymentStatus {
	switch raw {
	case "succeeded":
		return model.PaymentStatusSucceeded
	case "failed", "payment_failed":
		return model.PaymentStatusFailed
	case "canceled", "cancelled":
		return model.PaymentStatusCanceled
	case "processing":
		return model.PaymentStatusProcessing
	default:
		return model.PaymentStatusPending
	}
}

func toJSONMap(src map[string]string) datatypes.JSONMap {
	m := datatypes.JSONMap{}
	for k, v := range src {
		m[k] = v
	}
	return m
}
