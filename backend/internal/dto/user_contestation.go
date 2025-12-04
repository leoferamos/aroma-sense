package dto

import (
	"time"

	"github.com/leoferamos/aroma-sense/internal/model"
)

type UserContestationResponse struct {
	ID          uint       `json:"id"`
	UserID      uint       `json:"user_id"`
	Reason      string     `json:"reason"`
	Status      string     `json:"status"`
	RequestedAt time.Time  `json:"requested_at"`
	ReviewedAt  *time.Time `json:"reviewed_at,omitempty"`
	ReviewedBy  *uint      `json:"reviewed_by,omitempty"`
	ReviewNotes *string    `json:"review_notes,omitempty"`
}

func UserContestationResponseFromModel(m *model.UserContestation) UserContestationResponse {
	return UserContestationResponse{
		ID:          m.ID,
		UserID:      m.UserID,
		Reason:      m.Reason,
		Status:      m.Status,
		RequestedAt: m.RequestedAt,
		ReviewedAt:  m.ReviewedAt,
		ReviewedBy:  m.ReviewedBy,
		ReviewNotes: m.ReviewNotes,
	}
}
