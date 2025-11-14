package service

import (
	"context"
	"errors"

	"github.com/leoferamos/aroma-sense/internal/dto"
	"github.com/leoferamos/aroma-sense/internal/model"
	"github.com/leoferamos/aroma-sense/internal/repository"
	"github.com/leoferamos/aroma-sense/internal/validation"
)

// Sentinel errors for shipping flow to avoid string matching in handlers.
var (
	ErrUnauthorized        = errors.New("unauthorized")
	ErrCartEmpty           = errors.New("cart is empty")
	ErrInvalidPostalCode   = errors.New("invalid destination postal code")
	ErrOriginNotConfigured = errors.New("shipping origin not configured")
	ErrProviderUnavailable = errors.New("shipping provider not configured")
	ErrNoOptions           = errors.New("no shipping options available")
)

// ShippingProvider defines the interface for shipping quote providers.
type ShippingProvider interface {
	GetQuotes(ctx context.Context, userID string, originCEP, destCEP string, parcels []model.Parcel, insuredValue float64) ([]dto.ShippingOption, error)
}

// ShippingService exposes shipping quotations based on the user's cart and destination.
type ShippingService interface {
	CalculateOptions(ctx context.Context, userID string, postalCode string) ([]dto.ShippingOption, error)
}

type shippingService struct {
	cartRepo  repository.CartRepository
	provider  ShippingProvider
	originCEP string
}

// NewShippingService constructs a ShippingService.
func NewShippingService(cartRepo repository.CartRepository, provider ShippingProvider, originCEP string) ShippingService {
	return &shippingService{cartRepo: cartRepo, provider: provider, originCEP: originCEP}
}

func (s *shippingService) CalculateOptions(ctx context.Context, userID string, postalCode string) ([]dto.ShippingOption, error) {
	if userID == "" {
		return nil, ErrUnauthorized
	}
	if s.originCEP == "" {
		return nil, ErrOriginNotConfigured
	}
	destCEP := validation.NormalizeCEP(postalCode)
	if len(destCEP) < 5 {
		return nil, ErrInvalidPostalCode
	}

	cart, err := s.cartRepo.FindByUserID(userID)
	if err != nil || cart == nil || len(cart.Items) == 0 {
		return nil, ErrCartEmpty
	}

	// Aggregate weight and derive a single parcel.
	var totalWeightKg float64
	var insuredValue float64
	for _, it := range cart.Items {
		if it.Product != nil {
			w := it.Product.Weight
			if w > 50 {
				w = w / 1000.0
			}
			totalWeightKg += w * float64(it.Quantity)
			insuredValue += it.Price * float64(it.Quantity)
		} else {
			insuredValue += it.Price * float64(it.Quantity)
		}
	}
	if totalWeightKg <= 0 {
		totalWeightKg = 0.3
	}
	parcel := model.Parcel{WeightKg: totalWeightKg, LengthCm: 20, WidthCm: 15, HeightCm: 10}

	if s.provider == nil {
		return nil, ErrProviderUnavailable
	}
	quotes, err := s.provider.GetQuotes(ctx, userID, s.originCEP, destCEP, []model.Parcel{parcel}, insuredValue)
	if err != nil {
		return nil, err
	}
	if len(quotes) == 0 {
		return nil, ErrNoOptions
	}
	return quotes, nil
}
