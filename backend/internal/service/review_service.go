package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/leoferamos/aroma-sense/internal/apperror"
	"github.com/leoferamos/aroma-sense/internal/model"
	"github.com/leoferamos/aroma-sense/internal/repository"
)

// ReviewService defines business logic for product reviews
type ReviewService interface {
	CanUserReview(ctx context.Context, user *model.User, productID uint) (bool, string, error)
	CanUserReviewBySlug(ctx context.Context, user *model.User, slug string) (bool, string, error)
	CreateReview(ctx context.Context, user *model.User, productID uint, rating int, comment string) (*model.Review, error)
	ListReviews(ctx context.Context, productID uint, page, perPage int) ([]model.Review, int, error)
	GetAverage(ctx context.Context, productID uint) (float64, int, error)
	DeleteOwnReview(ctx context.Context, reviewID string, userID string) error
}

type ratingCacheEntry struct {
	avg     float64
	count   int
	staleAt time.Time
}

type reviewService struct {
	reviews  repository.ReviewRepository
	orders   repository.OrderRepository
	products repository.ProductRepository

	mu    sync.RWMutex
	cache map[uint]ratingCacheEntry
	ttl   time.Duration
}

func NewReviewService(reviews repository.ReviewRepository, orders repository.OrderRepository, products repository.ProductRepository) ReviewService {
	return &reviewService{
		reviews:  reviews,
		orders:   orders,
		products: products,
		cache:    make(map[uint]ratingCacheEntry),
		ttl:      5 * time.Minute,
	}
}

// CanUserReview checks if a user is eligible to review a product
func (s *reviewService) CanUserReview(ctx context.Context, user *model.User, productID uint) (bool, string, error) {
	// Must be authenticated
	if user == nil || user.PublicID == "" {
		return false, "unauthenticated", nil
	}
	// Must have display name
	if user.DisplayName == nil || strings.TrimSpace(*user.DisplayName) == "" {
		return false, "profile_incomplete", nil
	}
	// Must have a delivered order with the product
	delivered, err := s.orders.HasUserDeliveredOrderWithProduct(user.PublicID, productID)
	if err != nil {
		return false, "internal_error", apperror.NewDomain(fmt.Errorf("failed to verify delivered orders: %w", err), "internal_error", "internal error")
	}
	if !delivered {
		return false, "not_delivered", nil
	}
	// Must not have an existing review
	exists, err := s.reviews.ExistsByProductAndUser(ctx, productID, user.PublicID)
	if err != nil {
		return false, "internal_error", apperror.NewDomain(fmt.Errorf("failed to check existing review: %w", err), "internal_error", "internal error")
	}
	if exists {
		return false, "already_reviewed", nil
	}
	return true, "", nil
}

// CanUserReviewBySlug checks if a user can review a product identified by slug
func (s *reviewService) CanUserReviewBySlug(ctx context.Context, user *model.User, slug string) (bool, string, error) {
	prod, err := s.products.FindBySlug(slug)
	if err != nil {
		return false, "product_not_found", apperror.NewDomain(fmt.Errorf("failed to resolve product slug: %w", err), "product_not_found", "product not found")
	}
	return s.CanUserReview(ctx, user, prod.ID)
}

// CreateReview creates a new product review
func (s *reviewService) CreateReview(ctx context.Context, user *model.User, productID uint, rating int, comment string) (*model.Review, error) {
	if user == nil || user.PublicID == "" {
		return nil, apperror.NewCodeMessage("unauthenticated", "authentication required")
	}
	// Validate rating bounds
	if rating < 1 || rating > 5 {
		return nil, apperror.NewDomain(fmt.Errorf("invalid rating: %d", rating), "invalid_rating", "invalid rating")
	}
	// Validate comment length (<=500)
	if len(comment) > 500 {
		return nil, apperror.NewDomain(fmt.Errorf("comment too long: %d", len(comment)), "comment_too_long", "comment too long")
	}
	// Require display name
	if user.DisplayName == nil || strings.TrimSpace(*user.DisplayName) == "" {
		return nil, apperror.NewCodeMessage("profile_incomplete", "profile incomplete")
	}
	// Check delivered order
	delivered, err := s.orders.HasUserDeliveredOrderWithProduct(user.PublicID, productID)
	if err != nil {
		return nil, apperror.NewDomain(fmt.Errorf("failed to verify delivered orders: %w", err), "internal_error", "internal error")
	}
	if !delivered {
		return nil, apperror.NewCodeMessage("not_delivered", "product not delivered")
	}
	// Prevent duplicate review
	exists, err := s.reviews.ExistsByProductAndUser(ctx, productID, user.PublicID)
	if err != nil {
		return nil, apperror.NewDomain(fmt.Errorf("failed to check existing review: %w", err), "internal_error", "internal error")
	}
	if exists {
		return nil, apperror.NewCodeMessage("already_reviewed", "already reviewed")
	}

	// Verify product exists
	if _, err := s.products.FindByID(productID); err != nil {
		return nil, apperror.NewDomain(fmt.Errorf("product not found: %v", err), "product_not_found", "product not found")
	}

	rv := &model.Review{
		ProductID: productID,
		UserID:    user.PublicID,
		Rating:    rating,
		Comment:   strings.TrimSpace(comment),
		Status:    model.ReviewStatusPublished,
	}
	if err := s.reviews.CreateReview(ctx, rv); err != nil {
		return nil, apperror.NewDomain(fmt.Errorf("failed to create review: %w", err), "internal_error", "internal error")
	}

	// Invalidate cache
	s.mu.Lock()
	delete(s.cache, productID)
	s.mu.Unlock()

	return rv, nil
}

// ListReviews lists reviews for a product with pagination
func (s *reviewService) ListReviews(ctx context.Context, productID uint, page, perPage int) ([]model.Review, int, error) {
	if page < 1 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}
	const maxPerPage = 50
	if perPage > maxPerPage {
		perPage = maxPerPage
	}
	offset := (page - 1) * perPage
	return s.reviews.ListByProduct(ctx, productID, perPage, offset)
}

// GetAverage returns the average rating and count for a product.
func (s *reviewService) GetAverage(ctx context.Context, productID uint) (float64, int, error) {
	s.mu.RLock()
	entry, ok := s.cache[productID]
	isFresh := ok && time.Now().Before(entry.staleAt)
	s.mu.RUnlock()
	if isFresh {
		return entry.avg, entry.count, nil
	}
	s.mu.Lock()
	entry, ok = s.cache[productID]
	isFresh = ok && time.Now().Before(entry.staleAt)
	if isFresh {
		avg, count := entry.avg, entry.count
		s.mu.Unlock()
		return avg, count, nil
	}
	// Cache is still stale, query database
	avg, count, err := s.reviews.AverageRating(ctx, productID)
	if err != nil {
		s.mu.Unlock()
		return 0, 0, apperror.NewDomain(fmt.Errorf("failed to compute average rating: %w", err), "internal_error", "internal error")
	}
	s.cache[productID] = ratingCacheEntry{avg: avg, count: count, staleAt: time.Now().Add(s.ttl)}
	s.mu.Unlock()
	return avg, count, nil
}

// DeleteOwnReview allows a user to soft delete their own review
func (s *reviewService) DeleteOwnReview(ctx context.Context, reviewID string, userID string) error {
	// Try to get productID for precise cache invalidation
	var productID uint
	if pid, err := s.reviews.GetProductIDForUserReview(ctx, reviewID, userID); err == nil {
		productID = pid
	}

	if err := s.reviews.SoftDeleteReview(ctx, reviewID, userID); err != nil {
		if errors.Is(err, repository.ErrReviewNotFound) {
			return apperror.NewCodeMessage("review_not_found", "review not found")
		}
		return apperror.NewDomain(fmt.Errorf("failed to delete review: %w", err), "internal_error", "internal error")
	}

	if productID != 0 {
		s.mu.Lock()
		delete(s.cache, productID)
		s.mu.Unlock()
	}
	return nil
}
