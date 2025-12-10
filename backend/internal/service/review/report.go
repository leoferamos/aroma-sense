package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/leoferamos/aroma-sense/internal/apperror"
	"github.com/leoferamos/aroma-sense/internal/model"
	"github.com/leoferamos/aroma-sense/internal/repository"
	adminservice "github.com/leoferamos/aroma-sense/internal/service/admin"
	"gorm.io/gorm"
)

type ReviewReportService interface {
	Report(ctx context.Context, reviewID string, reporterID string, category string, reason string) error
	List(ctx context.Context, status string, limit, offset int) ([]model.ReviewReport, int64, error)
	Resolve(ctx context.Context, reportID string, action string, deactivateUser bool, adminPublicID string, suspensionUntil *time.Time, notes *string) error
}

type reviewReportService struct {
	reports   repository.ReviewReportRepository
	reviews   repository.ReviewRepository
	users     repository.UserRepository
	adminUser adminservice.AdminUserService
}

func NewReviewReportService(reports repository.ReviewReportRepository, reviews repository.ReviewRepository, users repository.UserRepository, adminUser adminservice.AdminUserService) ReviewReportService {
	return &reviewReportService{reports: reports, reviews: reviews, users: users, adminUser: adminUser}
}

var allowedCategories = map[string]struct{}{
	"offensive": {},
	"spam":      {},
	"improper":  {},
	"other":     {},
}

var allowedStatuses = map[string]struct{}{
	"pending":  {},
	"accepted": {},
	"rejected": {},
}

var allowedActions = map[string]string{
	"accept": "accepted",
	"reject": "rejected",
}

// Report allows a user to report a review for moderation
func (s *reviewReportService) Report(ctx context.Context, reviewID string, reporterID string, category string, reason string) error {
	category = strings.ToLower(strings.TrimSpace(category))
	reason = strings.TrimSpace(reason)

	if _, ok := allowedCategories[category]; !ok {
		return apperror.NewCodeMessage("invalid_category", "invalid category")
	}
	if len(reason) > 500 {
		return apperror.NewCodeMessage("reason_too_long", "reason too long")
	}

	review, err := s.reviews.FindByID(ctx, reviewID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return apperror.NewCodeMessage("review_not_found", "review not found")
		}
		return apperror.NewDomain(fmt.Errorf("find review: %w", err), "internal_error", "internal error")
	}

	if review.UserID == reporterID {
		return apperror.NewCodeMessage("cannot_report_own_review", "cannot report own review")
	}

	exists, err := s.reports.ExistsByReviewAndReporter(ctx, reviewID, reporterID)
	if err != nil {
		return apperror.NewDomain(fmt.Errorf("check duplicate report: %w", err), "internal_error", "internal error")
	}
	if exists {
		return apperror.NewCodeMessage("already_reported", "already reported")
	}

	report := &model.ReviewReport{
		ReviewID:       reviewID,
		ReportedBy:     reporterID,
		ReasonCategory: category,
		ReasonText:     reason,
		Status:         "pending",
	}
	if err := s.reports.Create(ctx, report); err != nil {
		return apperror.NewDomain(fmt.Errorf("create report: %w", err), "internal_error", "internal error")
	}

	if err := s.reports.IncrementReviewReportsCount(ctx, reviewID); err != nil {
		return apperror.NewDomain(fmt.Errorf("increment reports count: %w", err), "internal_error", "internal error")
	}

	return nil
}

func (s *reviewReportService) List(ctx context.Context, status string, limit, offset int) ([]model.ReviewReport, int64, error) {
	status = strings.ToLower(strings.TrimSpace(status))
	if status == "" {
		status = "pending"
	}
	if _, ok := allowedStatuses[status]; !ok {
		return nil, 0, apperror.NewCodeMessage("invalid_status", "invalid status")
	}

	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	reports, total, err := s.reports.ListByStatus(ctx, status, limit, offset)
	if err != nil {
		return nil, 0, apperror.NewDomain(fmt.Errorf("list reports: %w", err), "internal_error", "internal error")
	}

	return reports, total, nil
}

func (s *reviewReportService) Resolve(ctx context.Context, reportID string, action string, deactivateUser bool, adminPublicID string, suspensionUntil *time.Time, notes *string) error {
	status, ok := allowedActions[strings.ToLower(strings.TrimSpace(action))]
	if !ok {
		return apperror.NewCodeMessage("invalid_action", "invalid action")
	}

	report, err := s.reports.GetByID(ctx, reportID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return apperror.NewCodeMessage("report_not_found", "report not found")
		}
		return apperror.NewDomain(fmt.Errorf("get report: %w", err), "internal_error", "internal error")
	}

	if report.Status != "pending" {
		return apperror.NewCodeMessage("report_already_resolved", "report already resolved")
	}

	// Accept path hides the review
	if status == "accepted" {
		if err := s.reviews.UpdateStatus(ctx, report.ReviewID, model.ReviewStatusHidden); err != nil {
			if err == repository.ErrReviewNotFound {
				return apperror.NewCodeMessage("review_not_found", "review not found")
			}
			return apperror.NewDomain(fmt.Errorf("hide review: %w", err), "internal_error", "internal error")
		}
	}

	if err := s.reports.UpdateStatus(ctx, reportID, status); err != nil {
		return apperror.NewDomain(fmt.Errorf("update report status: %w", err), "internal_error", "internal error")
	}

	if deactivateUser && s.adminUser != nil {
		// Review.UserID stores public_id; fetch user to get numeric ID
		reviewAuthor, err := s.users.FindByPublicID(report.Review.UserID)
		if err != nil {
			return apperror.NewDomain(fmt.Errorf("get review author: %w", err), "internal_error", "internal error")
		}
		reason := "reported_review"
		notesVal := ""
		if notes != nil {
			notesVal = *notes
		}
		if err := s.adminUser.DeactivateUser(reviewAuthor.ID, adminPublicID, reason, notesVal, suspensionUntil); err != nil {
			return apperror.NewDomain(fmt.Errorf("deactivate user: %w", err), "internal_error", "internal error")
		}
	}

	return nil
}
