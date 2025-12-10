package repository

import (
	"context"

	"github.com/leoferamos/aroma-sense/internal/model"
	"gorm.io/gorm"
)

type ReviewReportRepository interface {
	Create(ctx context.Context, report *model.ReviewReport) error
	ExistsByReviewAndReporter(ctx context.Context, reviewID, reporterID string) (bool, error)
	IncrementReviewReportsCount(ctx context.Context, reviewID string) error
	ListByStatus(ctx context.Context, status string, limit, offset int) ([]model.ReviewReport, int64, error)
	GetByID(ctx context.Context, id string) (*model.ReviewReport, error)
	UpdateStatus(ctx context.Context, id string, status string) error
}

type reviewReportRepository struct {
	db *gorm.DB
}

func NewReviewReportRepository(db *gorm.DB) ReviewReportRepository {
	return &reviewReportRepository{db: db}
}

// Create inserts a new review report into the database
func (r *reviewReportRepository) Create(ctx context.Context, report *model.ReviewReport) error {
	return r.db.WithContext(ctx).Create(report).Error
}

// ExistsByReviewAndReporter checks if a report already exists for the given review and reporter
func (r *reviewReportRepository) ExistsByReviewAndReporter(ctx context.Context, reviewID, reporterID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM review_reports WHERE review_id = ? AND reported_by = ?)`
	if err := r.db.WithContext(ctx).Raw(query, reviewID, reporterID).Scan(&exists).Error; err != nil {
		return false, err
	}
	return exists, nil
}

// IncrementReviewReportsCount increments the reports_count for a review
func (r *reviewReportRepository) IncrementReviewReportsCount(ctx context.Context, reviewID string) error {
	return r.db.WithContext(ctx).Model(&model.Review{}).
		Where("id = ?", reviewID).
		UpdateColumn("reports_count", gorm.Expr("reports_count + 1")).Error
}

// GetByID fetches a review report by ID
func (r *reviewReportRepository) GetByID(ctx context.Context, id string) (*model.ReviewReport, error) {
	var report model.ReviewReport
	if err := r.db.WithContext(ctx).Preload("Review").Preload("Reporter").First(&report, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &report, nil
}

// UpdateStatus updates the status of a report
func (r *reviewReportRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	return r.db.WithContext(ctx).Model(&model.ReviewReport{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// ListByStatus returns reports filtered by status with pagination and total count
func (r *reviewReportRepository) ListByStatus(ctx context.Context, status string, limit, offset int) ([]model.ReviewReport, int64, error) {
	var total int64
	query := r.db.WithContext(ctx).Model(&model.ReviewReport{}).Where("status = ?", status)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var reports []model.ReviewReport
	if err := query.Preload("Review").Preload("Reporter").Order("created_at DESC").Limit(limit).Offset(offset).Find(&reports).Error; err != nil {
		return nil, 0, err
	}

	return reports, total, nil
}
