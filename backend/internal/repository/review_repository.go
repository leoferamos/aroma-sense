package repository

import (
	"context"
	"errors"
	"time"

	"github.com/leoferamos/aroma-sense/internal/model"
	"gorm.io/gorm"
)

var ErrReviewNotFound = errors.New("review not found")

type ReviewRepository interface {
	CreateReview(ctx context.Context, review *model.Review) error
	ListByProduct(ctx context.Context, productID uint, limit, offset int) ([]model.Review, int, error)
	AverageRating(ctx context.Context, productID uint) (float64, int, error)
	ExistsByProductAndUser(ctx context.Context, productID uint, userID string) (bool, error)
	GetProductIDForUserReview(ctx context.Context, reviewID string, userID string) (uint, error)
	SoftDeleteReview(ctx context.Context, reviewID string, userID string) error
	FindByID(ctx context.Context, reviewID string) (*model.Review, error)
	UpdateStatus(ctx context.Context, reviewID string, status model.ReviewStatus) error
}

type reviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) ReviewRepository {
	return &reviewRepository{db: db}
}

// CreateReview creates a new review record
func (r *reviewRepository) CreateReview(ctx context.Context, review *model.Review) error {
	return r.db.WithContext(ctx).Create(review).Error
}

// ListByProduct returns paginated list of published reviews for a product
func (r *reviewRepository) ListByProduct(ctx context.Context, productID uint, limit, offset int) ([]model.Review, int, error) {
	var reviews []model.Review
	q := r.db.WithContext(ctx).Model(&model.Review{}).
		Where("product_id = ? AND status = ? AND deleted_at IS NULL", productID, model.ReviewStatusPublished)

	// Count total
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	if err := q.Order("created_at DESC").
		Preload("User", func(db *gorm.DB) *gorm.DB { return db.Select("public_id", "display_name") }).
		Offset(offset).Limit(limit).
		Find(&reviews).Error; err != nil {
		return nil, 0, err
	}

	return reviews, int(total), nil
}

// AverageRating computes the average rating and total count of published reviews for a product
func (r *reviewRepository) AverageRating(ctx context.Context, productID uint) (float64, int, error) {
	type agg struct {
		Avg   *float64
		Count int
	}
	var a agg
	raw := `SELECT COALESCE(AVG(rating),0) AS avg, COUNT(*) AS count FROM reviews WHERE product_id = ? AND status = 'published' AND deleted_at IS NULL`
	if err := r.db.WithContext(ctx).Raw(raw, productID).Scan(&a).Error; err != nil {
		return 0, 0, err
	}
	if a.Avg == nil {
		zero := 0.0
		a.Avg = &zero
	}
	return *a.Avg, a.Count, nil
}

// ExistsByProductAndUser checks if a user has already reviewed a product
func (r *reviewRepository) ExistsByProductAndUser(ctx context.Context, productID uint, userID string) (bool, error) {
	var exists bool
	raw := `SELECT EXISTS(SELECT 1 FROM reviews WHERE product_id = ? AND user_id = ? AND deleted_at IS NULL)`
	if err := r.db.WithContext(ctx).Raw(raw, productID, userID).Scan(&exists).Error; err != nil {
		return false, err
	}
	return exists, nil
}

// SoftDeleteReview marks a review as deleted by setting deleted_at
func (r *reviewRepository) SoftDeleteReview(ctx context.Context, reviewID string, userID string) error {
	now := time.Now()
	result := r.db.WithContext(ctx).Model(&model.Review{}).
		Where("id = ? AND user_id = ? AND deleted_at IS NULL", reviewID, userID).
		Updates(map[string]interface{}{"deleted_at": now})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrReviewNotFound
	}
	return nil
}

// GetProductIDForUserReview retrieves the product ID for a given review ID and user ID
func (r *reviewRepository) GetProductIDForUserReview(ctx context.Context, reviewID string, userID string) (uint, error) {
	var res struct {
		ProductID uint
	}
	if err := r.db.WithContext(ctx).
		Model(&model.Review{}).
		Select("product_id").
		Where("id = ? AND user_id = ? AND deleted_at IS NULL", reviewID, userID).
		First(&res).Error; err != nil {
		return 0, err
	}
	return res.ProductID, nil
}

// FindByID returns a review by ID when not soft-deleted
func (r *reviewRepository) FindByID(ctx context.Context, reviewID string) (*model.Review, error) {
	var review model.Review
	if err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", reviewID).
		First(&review).Error; err != nil {
		return nil, err
	}
	return &review, nil
}

// UpdateStatus updates the status of a review (e.g., published/hidden)
func (r *reviewRepository) UpdateStatus(ctx context.Context, reviewID string, status model.ReviewStatus) error {
	result := r.db.WithContext(ctx).
		Model(&model.Review{}).
		Where("id = ? AND deleted_at IS NULL", reviewID).
		Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrReviewNotFound
	}
	return nil
}
