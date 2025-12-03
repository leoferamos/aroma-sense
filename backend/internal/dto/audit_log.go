package dto

import (
	"time"

	"github.com/leoferamos/aroma-sense/internal/model"
)

// UserBasicResponse represents basic user information for audit logs
type UserBasicResponse struct {
	ID          uint   `json:"id"`
	PublicID    string `json:"public_id"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name,omitempty"`
	Role        string `json:"role"`
}

// AuditLogResponse represents audit log data for API responses
type AuditLogResponse struct {
	ID         uint                   `json:"id"`
	PublicID   string                 `json:"public_id"`
	UserID     *uint                  `json:"user_id,omitempty"`
	ActorID    *uint                  `json:"actor_id,omitempty"`
	ActorType  string                 `json:"actor_type"`
	Action     model.AuditAction      `json:"action"`
	Resource   string                 `json:"resource"`
	ResourceID *string                `json:"resource_id,omitempty"`
	Details    map[string]interface{} `json:"details,omitempty"`
	OldValues  map[string]interface{} `json:"old_values,omitempty"`
	NewValues  map[string]interface{} `json:"new_values,omitempty"`
	Timestamp  time.Time              `json:"timestamp"`
	Compliance string                 `json:"compliance,omitempty"`
	Severity   string                 `json:"severity"`
	CreatedAt  time.Time              `json:"created_at"`
	User       *UserBasicResponse     `json:"user,omitempty"`
	Actor      *UserBasicResponse     `json:"actor,omitempty"`
}

// AuditLogListResponse represents paginated audit log list
type AuditLogListResponse struct {
	AuditLogs []AuditLogResponse `json:"audit_logs"`
	Total     int64              `json:"total"`
	Limit     int                `json:"limit"`
	Offset    int                `json:"offset"`
}

// AuditLogFilterRequest represents filters for audit log queries
type AuditLogFilterRequest struct {
	UserID     *uint      `json:"user_id,omitempty" form:"user_id"`
	ActorID    *uint      `json:"actor_id,omitempty" form:"actor_id"`
	Action     *string    `json:"action,omitempty" form:"action" validate:"omitempty,oneof=user_login user_logout user_created user_updated user_deactivated user_reactivated user_deleted user_exported role_changed data_accessed deletion_requested deletion_confirmed deletion_cancelled data_anonymized"`
	Resource   *string    `json:"resource,omitempty" form:"resource"`
	ResourceID *string    `json:"resource_id,omitempty" form:"resource_id"`
	StartDate  *time.Time `json:"start_date,omitempty" form:"start_date"`
	EndDate    *time.Time `json:"end_date,omitempty" form:"end_date"`
	Severity   *string    `json:"severity,omitempty" form:"severity" validate:"omitempty,oneof=info warning error critical"`
	Limit      int        `json:"limit,omitempty" form:"limit" validate:"min=1,max=1000"`
	Offset     int        `json:"offset,omitempty" form:"offset" validate:"min=0"`
}

// AuditLogSummaryResponse represents audit log summary statistics
type AuditLogSummaryResponse struct {
	TotalActions  int64              `json:"total_actions"`
	ActionsByType map[string]int64   `json:"actions_by_type"`
	RecentActions []AuditLogResponse `json:"recent_actions"`
	UserActivity  map[string]int64   `json:"user_activity"`
	GeneratedAt   time.Time          `json:"generated_at"`
}

// AuditLogCreateRequest represents data for creating audit log entries
type AuditLogCreateRequest struct {
	UserID     *uint                  `json:"user_id,omitempty"`
	ActorID    *uint                  `json:"actor_id,omitempty"`
	Action     model.AuditAction      `json:"action" validate:"required"`
	Resource   string                 `json:"resource" validate:"required"`
	ResourceID *string                `json:"resource_id,omitempty"`
	Details    map[string]interface{} `json:"details,omitempty"`
	Severity   string                 `json:"severity,omitempty" validate:"omitempty,oneof=info warning error critical"`
}
