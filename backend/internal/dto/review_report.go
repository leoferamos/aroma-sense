package dto

import (
	"time"

	"github.com/leoferamos/aroma-sense/internal/model"
)

// ReviewReportRequest is the payload for reporting a review.
type ReviewReportRequest struct {
	Category string `json:"category" binding:"required"`
	Reason   string `json:"reason" binding:"omitempty,max=500"`
}

// ReviewReportAdminReview summarizes the reported review for admin list
type ReviewReportAdminReview struct {
	ID      string `json:"id"`
	Comment string `json:"comment"`
	Rating  int    `json:"rating"`
	UserID  string `json:"user_id"`
}

// ReviewReportAdminReporter summarizes the reporter user for admin list
type ReviewReportAdminReporter struct {
	PublicID    string  `json:"public_id"`
	DisplayName *string `json:"display_name"`
}

// ReviewReportAdminItem represents a review report in admin listings
type ReviewReportAdminItem struct {
	ID             string                     `json:"id"`
	ReviewID       string                     `json:"review_id"`
	ReportedBy     string                     `json:"reported_by"`
	ReasonCategory string                     `json:"reason_category"`
	ReasonText     string                     `json:"reason_text"`
	Status         string                     `json:"status"`
	CreatedAt      time.Time                  `json:"created_at"`
	Review         *ReviewReportAdminReview   `json:"review,omitempty"`
	Reporter       *ReviewReportAdminReporter `json:"reporter,omitempty"`
}

// ReviewReportAdminResponse wraps paginated admin list results
type ReviewReportAdminResponse struct {
	Items  []ReviewReportAdminItem `json:"items"`
	Total  int64                   `json:"total"`
	Limit  int                     `json:"limit"`
	Offset int                     `json:"offset"`
}

// ReviewReportResolveRequest is used by admin to resolve a report
type ReviewReportResolveRequest struct {
	Action          string  `json:"action" binding:"required,oneof=accept reject"`
	DeactivateUser  bool    `json:"deactivate_user"`
	SuspensionUntil *string `json:"suspension_until"`
	Notes           *string `json:"notes"`
}

// ReviewReportAdminItemFromModel converts model to admin DTO
func ReviewReportAdminItemFromModel(m *model.ReviewReport) ReviewReportAdminItem {
	item := ReviewReportAdminItem{
		ID:             m.ID,
		ReviewID:       m.ReviewID,
		ReportedBy:     m.ReportedBy,
		ReasonCategory: m.ReasonCategory,
		ReasonText:     m.ReasonText,
		Status:         m.Status,
		CreatedAt:      m.CreatedAt,
	}

	if m.Review != nil {
		item.Review = &ReviewReportAdminReview{
			ID:      m.Review.ID,
			Comment: m.Review.Comment,
			Rating:  m.Review.Rating,
			UserID:  m.Review.UserID,
		}
	}

	if m.Reporter != nil {
		item.Reporter = &ReviewReportAdminReporter{
			PublicID:    m.Reporter.PublicID,
			DisplayName: m.Reporter.DisplayName,
		}
	}

	return item
}
