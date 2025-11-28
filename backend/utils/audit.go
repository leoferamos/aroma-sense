package utils

import (
	"encoding/json"
	"log"
	"time"

	"github.com/leoferamos/aroma-sense/internal/dto"
	"github.com/leoferamos/aroma-sense/internal/model"
)

// ParseJSONField safely parses a JSON string into a map
func ParseJSONField(jsonStr string) map[string]interface{} {
	if jsonStr == "" {
		return nil
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		log.Printf("ERROR: Failed to parse JSON field in audit log: %v (field length: %d)", err, len(jsonStr))
		return nil
	}

	return result
}

// ConvertUserToBasicResponse converts model.User to dto.UserBasicResponse
func ConvertUserToBasicResponse(user *model.User) dto.UserBasicResponse {
	if user == nil {
		return dto.UserBasicResponse{}
	}

	displayName := ""
	if user.DisplayName != nil {
		displayName = *user.DisplayName
	}

	return dto.UserBasicResponse{
		ID:          user.ID,
		PublicID:    user.PublicID,
		Email:       MaskEmail(user.Email), // LGPD: Mask email in logs list view
		DisplayName: displayName,
		Role:        user.Role,
	}
}

// ConvertUserToBasicResponseDetailed converts model.User to dto.UserBasicResponse with full email
func ConvertUserToBasicResponseDetailed(user *model.User) dto.UserBasicResponse {
	if user == nil {
		return dto.UserBasicResponse{}
	}

	displayName := ""
	if user.DisplayName != nil {
		displayName = *user.DisplayName
	}

	return dto.UserBasicResponse{
		ID:          user.ID,
		PublicID:    user.PublicID,
		Email:       user.Email, // Full email for detailed operational view
		DisplayName: displayName,
		Role:        user.Role,
	}
}

// ConvertAuditLogToResponse converts model.AuditLog to dto.AuditLogResponse
func ConvertAuditLogToResponse(auditLog *model.AuditLog) dto.AuditLogResponse {
	response := dto.AuditLogResponse{
		ID:         auditLog.ID,
		PublicID:   auditLog.PublicID.String(),
		UserID:     auditLog.UserID,
		ActorID:    auditLog.ActorID,
		ActorType:  auditLog.ActorType,
		Action:     auditLog.Action,
		Resource:   auditLog.Resource,
		ResourceID: auditLog.ResourceID,
		Timestamp:  auditLog.Timestamp,
		Compliance: auditLog.Compliance,
		Severity:   auditLog.Severity,
		CreatedAt:  auditLog.CreatedAt,
	}

	// Parse JSON fields safely
	response.Details = ParseJSONField(auditLog.Details)
	response.OldValues = ParseJSONField(auditLog.OldValues)
	response.NewValues = ParseJSONField(auditLog.NewValues)

	if auditLog.User != nil {
		userResponse := ConvertUserToBasicResponse(auditLog.User)
		response.User = &userResponse
	}

	if auditLog.Actor != nil {
		actorResponse := ConvertUserToBasicResponse(auditLog.Actor)
		response.Actor = &actorResponse
	}

	return response
}

// ConvertAuditLogToResponseDetailed converts model.AuditLog to dto.AuditLogResponse with full emails
func ConvertAuditLogToResponseDetailed(auditLog *model.AuditLog) dto.AuditLogResponse {
	response := dto.AuditLogResponse{
		ID:         auditLog.ID,
		PublicID:   auditLog.PublicID.String(),
		UserID:     auditLog.UserID,
		ActorID:    auditLog.ActorID,
		ActorType:  auditLog.ActorType,
		Action:     auditLog.Action,
		Resource:   auditLog.Resource,
		ResourceID: auditLog.ResourceID,
		Timestamp:  auditLog.Timestamp,
		Compliance: auditLog.Compliance,
		Severity:   auditLog.Severity,
		CreatedAt:  auditLog.CreatedAt,
	}

	// Parse JSON fields safely
	response.Details = ParseJSONField(auditLog.Details)
	response.OldValues = ParseJSONField(auditLog.OldValues)
	response.NewValues = ParseJSONField(auditLog.NewValues)

	// Convert related entities with FULL emails for operational purposes
	if auditLog.User != nil {
		userResponse := ConvertUserToBasicResponseDetailed(auditLog.User)
		response.User = &userResponse
	}

	if auditLog.Actor != nil {
		actorResponse := ConvertUserToBasicResponseDetailed(auditLog.Actor)
		response.Actor = &actorResponse
	}

	return response
}

// ConvertAuditLogsToResponse converts multiple audit logs to responses
func ConvertAuditLogsToResponse(auditLogs []*model.AuditLog) []dto.AuditLogResponse {
	responses := make([]dto.AuditLogResponse, len(auditLogs))
	for i, auditLog := range auditLogs {
		responses[i] = ConvertAuditLogToResponse(auditLog)
	}
	return responses
}

// ConvertAuditLogsToResponseDetailed converts multiple audit logs to responses with full emails
func ConvertAuditLogsToResponseDetailed(auditLogs []*model.AuditLog) []dto.AuditLogResponse {
	responses := make([]dto.AuditLogResponse, len(auditLogs))
	for i, auditLog := range auditLogs {
		responses[i] = ConvertAuditLogToResponseDetailed(auditLog)
	}
	return responses
}

// ConvertAuditLogSummaryToResponse converts summary to response
func ConvertAuditLogSummaryToResponse(summary *model.AuditLogSummary) dto.AuditLogSummaryResponse {
	// Convert []model.AuditLog to []*model.AuditLog
	recentActions := make([]*model.AuditLog, len(summary.RecentActions))
	for i := range summary.RecentActions {
		recentActions[i] = &summary.RecentActions[i]
	}

	return dto.AuditLogSummaryResponse{
		TotalActions:  summary.TotalActions,
		ActionsByType: summary.ActionsByType,
		RecentActions: ConvertAuditLogsToResponse(recentActions),
		UserActivity:  summary.UserActivity,
		GeneratedAt:   time.Now(),
	}
}
