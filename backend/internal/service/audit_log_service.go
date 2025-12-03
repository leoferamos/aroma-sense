package service

import (
	"encoding/json"
	"time"

	"github.com/leoferamos/aroma-sense/internal/dto"
	"github.com/leoferamos/aroma-sense/internal/model"
	"github.com/leoferamos/aroma-sense/internal/repository"
	"github.com/leoferamos/aroma-sense/internal/utils"
)

// AuditLogService defines the interface for audit log business logic
type AuditLogService interface {
	LogUserAction(actorID *uint, userID *uint, action model.AuditAction, details map[string]interface{}) error
	LogUserUpdate(actorID uint, userID uint, oldUser, newUser *model.User) error
	LogAdminAction(adminID uint, userID uint, action model.AuditAction, details map[string]interface{}) error
	LogSystemAction(action model.AuditAction, resource, resourceID string, details map[string]interface{}) error
	LogDataAccess(userID uint, resource string, resourceID string) error
	LogDeletionAction(actorID *uint, userID uint, action model.AuditAction, details map[string]interface{}) error
	ListAuditLogs(filter *model.AuditLogFilter) ([]*model.AuditLog, int64, error)
	GetAuditLogByID(id uint) (*model.AuditLog, error)
	GetUserAuditLogs(userID uint, limit, offset int) ([]*model.AuditLog, int64, error)
	GetResourceAuditLogs(resource, resourceID string) ([]*model.AuditLog, error)
	GetAuditSummary(startDate, endDate *time.Time) (*model.AuditLogSummary, error)
	CleanupOldLogs(retentionDays int) error
	ConvertAuditLogToResponse(auditLog *model.AuditLog) dto.AuditLogResponse
	ConvertAuditLogToResponseDetailed(auditLog *model.AuditLog) dto.AuditLogResponse
	ConvertAuditLogsToResponse(auditLogs []*model.AuditLog) []dto.AuditLogResponse
	ConvertAuditLogSummaryToResponse(summary *model.AuditLogSummary) dto.AuditLogSummaryResponse
}

type auditLogService struct {
	repo repository.AuditLogRepository
}

// NewAuditLogService creates a new instance of AuditLogService
func NewAuditLogService(repo repository.AuditLogRepository) AuditLogService {
	return &auditLogService{repo: repo}
}

// LogUserAction logs a general user action
func (s *auditLogService) LogUserAction(actorID *uint, userID *uint, action model.AuditAction, details map[string]interface{}) error {
	return s.logAuditEntry(actorID, userID, action, "user", nil, details, nil, nil)
}

// LogUserUpdate logs user profile updates with before/after values
func (s *auditLogService) LogUserUpdate(actorID uint, userID uint, oldUser, newUser *model.User) error {
	oldValues := map[string]interface{}{
		"email":          utils.MaskEmail(oldUser.Email), // LGPD: Mask email in logs
		"display_name":   oldUser.DisplayName,
		"role":           oldUser.Role,
		"deactivated_at": oldUser.DeactivatedAt,
	}

	newValues := map[string]interface{}{
		"email":          utils.MaskEmail(newUser.Email), // LGPD: Mask email in logs
		"display_name":   newUser.DisplayName,
		"role":           newUser.Role,
		"deactivated_at": newUser.DeactivatedAt,
	}

	details := map[string]interface{}{
		"action_type": "user_update",
		"actor_type":  "admin",
	}

	return s.logAuditEntry(&actorID, &userID, model.AuditActionUserUpdated, "user", nil, details, oldValues, newValues)
}

// LogAdminAction logs administrative actions on users
func (s *auditLogService) LogAdminAction(adminID uint, userID uint, action model.AuditAction, details map[string]interface{}) error {
	details["admin_action"] = true
	details["actor_type"] = "admin"

	return s.logAuditEntry(&adminID, &userID, action, "user", nil, details, nil, nil)
}

// LogSystemAction logs system-level actions
func (s *auditLogService) LogSystemAction(action model.AuditAction, resource, resourceID string, details map[string]interface{}) error {
	return s.logAuditEntry(nil, nil, action, resource, &resourceID, details, nil, nil)
}

// LogDataAccess logs when user data is accessed (LGPD compliance)
func (s *auditLogService) LogDataAccess(userID uint, resource string, resourceID string) error {
	details := map[string]interface{}{
		"access_type": "data_access",
		"purpose":     "user_request",
	}

	return s.logAuditEntry(nil, &userID, model.AuditActionDataAccessed, resource, &resourceID, details, nil, nil)
}

// LogDeletionAction logs account deletion related actions
func (s *auditLogService) LogDeletionAction(actorID *uint, userID uint, action model.AuditAction, details map[string]interface{}) error {
	details["deletion_action"] = true
	details["lgpd_compliant"] = true

	return s.logAuditEntry(actorID, &userID, action, "user", nil, details, nil, nil)
}

// logAuditEntry is a helper method to create audit log entries
func (s *auditLogService) logAuditEntry(actorID, userID *uint, action model.AuditAction, resource string, resourceID *string, details, oldValues, newValues map[string]interface{}) error {
	// Convert maps to JSON strings
	detailsJSON := "{}"
	if details != nil {
		if jsonBytes, err := json.Marshal(details); err == nil {
			detailsJSON = string(jsonBytes)
		}
	}

	oldValuesJSON := ""
	if oldValues != nil {
		if jsonBytes, err := json.Marshal(oldValues); err == nil {
			oldValuesJSON = string(jsonBytes)
		}
	}

	newValuesJSON := ""
	if newValues != nil {
		if jsonBytes, err := json.Marshal(newValues); err == nil {
			newValuesJSON = string(jsonBytes)
		}
	}

	auditLog := &model.AuditLog{
		UserID:     userID,
		ActorID:    actorID,
		ActorType:  "user",
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		Details:    detailsJSON,
		OldValues:  oldValuesJSON,
		NewValues:  newValuesJSON,
		Timestamp:  time.Now(),
		Compliance: "LGPD",
		Severity:   "info",
	}

	if actorID == nil {
		auditLog.ActorType = "system"
	}

	return s.repo.Create(auditLog)
}

// ListAuditLogs returns paginated audit logs with filters
func (s *auditLogService) ListAuditLogs(filter *model.AuditLogFilter) ([]*model.AuditLog, int64, error) {
	return s.repo.List(filter)
}

// GetAuditLogByID retrieves a specific audit log
func (s *auditLogService) GetAuditLogByID(id uint) (*model.AuditLog, error) {
	return s.repo.GetByID(id)
}

// GetUserAuditLogs retrieves audit logs for a specific user
func (s *auditLogService) GetUserAuditLogs(userID uint, limit, offset int) ([]*model.AuditLog, int64, error) {
	return s.repo.GetByUserID(userID, limit, offset)
}

// GetResourceAuditLogs retrieves audit logs for a specific resource
func (s *auditLogService) GetResourceAuditLogs(resource, resourceID string) ([]*model.AuditLog, error) {
	return s.repo.GetByResource(resource, resourceID)
}

// GetAuditSummary generates audit statistics
func (s *auditLogService) GetAuditSummary(startDate, endDate *time.Time) (*model.AuditLogSummary, error) {
	return s.repo.GetSummary(startDate, endDate)
}

// CleanupOldLogs removes old audit logs based on retention policy
func (s *auditLogService) CleanupOldLogs(retentionDays int) error {
	cutoffDate := time.Now().AddDate(0, 0, -retentionDays)
	return s.repo.DeleteOldLogs(cutoffDate)
}

// ConvertAuditLogToResponse converts model.AuditLog to dto.AuditLogResponse
func (s *auditLogService) ConvertAuditLogToResponse(auditLog *model.AuditLog) dto.AuditLogResponse {
	return utils.ConvertAuditLogToResponse(auditLog)
}

// ConvertAuditLogToResponseDetailed converts model.AuditLog to dto.AuditLogResponse with full emails
func (s *auditLogService) ConvertAuditLogToResponseDetailed(auditLog *model.AuditLog) dto.AuditLogResponse {
	return utils.ConvertAuditLogToResponseDetailed(auditLog)
}

// ConvertAuditLogsToResponse converts multiple audit logs to responses
func (s *auditLogService) ConvertAuditLogsToResponse(auditLogs []*model.AuditLog) []dto.AuditLogResponse {
	return utils.ConvertAuditLogsToResponse(auditLogs)
}

// ConvertAuditLogSummaryToResponse converts summary to response
func (s *auditLogService) ConvertAuditLogSummaryToResponse(summary *model.AuditLogSummary) dto.AuditLogSummaryResponse {
	return utils.ConvertAuditLogSummaryToResponse(summary)
}
