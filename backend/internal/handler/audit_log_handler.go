package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/dto"
	"github.com/leoferamos/aroma-sense/internal/model"
	"github.com/leoferamos/aroma-sense/internal/service"
)

type AuditLogHandler struct {
	auditLogService service.AuditLogService
}

// NewAuditLogHandler creates a new audit log handler
func NewAuditLogHandler(auditLogService service.AuditLogService) *AuditLogHandler {
	return &AuditLogHandler{
		auditLogService: auditLogService,
	}
}

// ListAuditLogs retrieves a paginated list of audit logs with filtering options
func (h *AuditLogHandler) ListAuditLogs(c *gin.Context) {
	// Parse query parameters
	filter := &dto.AuditLogFilterRequest{}

	if userIDStr := c.Query("user_id"); userIDStr != "" {
		if userID, err := strconv.ParseUint(userIDStr, 10, 32); err == nil {
			userIDUint := uint(userID)
			filter.UserID = &userIDUint
		}
	}

	if actorIDStr := c.Query("actor_id"); actorIDStr != "" {
		if actorID, err := strconv.ParseUint(actorIDStr, 10, 32); err == nil {
			actorIDUint := uint(actorID)
			filter.ActorID = &actorIDUint
		}
	}

	if action := c.Query("action"); action != "" {
		filter.Action = &action
	}

	if resource := c.Query("resource"); resource != "" {
		filter.Resource = &resource
	}

	if resourceID := c.Query("resource_id"); resourceID != "" {
		filter.ResourceID = &resourceID
	}

	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if startDate, err := time.Parse(time.RFC3339, startDateStr); err == nil {
			filter.StartDate = &startDate
		}
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if endDate, err := time.Parse(time.RFC3339, endDateStr); err == nil {
			filter.EndDate = &endDate
		}
	}

	if severity := c.Query("severity"); severity != "" {
		filter.Severity = &severity
	}

	// Pagination
	limit := 50
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 1000 {
			limit = l
		}
	}
	filter.Limit = limit

	offset := 0
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}
	filter.Offset = offset

	// Convert to model filter
	modelFilter := &model.AuditLogFilter{
		UserID:     filter.UserID,
		ActorID:    filter.ActorID,
		Resource:   filter.Resource,
		ResourceID: filter.ResourceID,
		StartDate:  filter.StartDate,
		EndDate:    filter.EndDate,
		Severity:   filter.Severity,
		Limit:      filter.Limit,
		Offset:     filter.Offset,
	}

	if filter.Action != nil {
		action := model.AuditAction(*filter.Action)
		modelFilter.Action = &action
	}

	// Get audit logs
	auditLogs, total, err := h.auditLogService.ListAuditLogs(modelFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to retrieve audit logs"})
		return
	}

	// Convert to response
	response := dto.AuditLogListResponse{
		AuditLogs: h.auditLogService.ConvertAuditLogsToResponse(auditLogs),
		Total:     total,
		Limit:     filter.Limit,
		Offset:    filter.Offset,
	}

	c.JSON(http.StatusOK, response)
}

// GetAuditLog retrieves a specific audit log entry with masked emails for general viewing
func (h *AuditLogHandler) GetAuditLog(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid audit log ID"})
		return
	}

	auditLog, err := h.auditLogService.GetAuditLogByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "Audit log not found"})
		return
	}

	response := h.auditLogService.ConvertAuditLogToResponse(auditLog)
	c.JSON(http.StatusOK, response)
}

// GetAuditLogDetailed retrieves a specific audit log entry with full emails for operational purposes
func (h *AuditLogHandler) GetAuditLogDetailed(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid audit log ID"})
		return
	}

	auditLog, err := h.auditLogService.GetAuditLogByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "Audit log not found"})
		return
	}

	response := h.auditLogService.ConvertAuditLogToResponseDetailed(auditLog)
	c.JSON(http.StatusOK, response)
}

// GetUserAuditLogs retrieves paginated audit logs for a specific user with masked emails
func (h *AuditLogHandler) GetUserAuditLogs(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid user ID"})
		return
	}

	limit := 50
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	offset := 0
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	auditLogs, total, err := h.auditLogService.GetUserAuditLogs(uint(userID), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to retrieve user audit logs"})
		return
	}

	response := dto.AuditLogListResponse{
		AuditLogs: h.auditLogService.ConvertAuditLogsToResponse(auditLogs),
		Total:     total,
		Limit:     limit,
		Offset:    offset,
	}

	c.JSON(http.StatusOK, response)
}

// GetAuditSummary returns audit logs summary and statistics within a date range
func (h *AuditLogHandler) GetAuditSummary(c *gin.Context) {
	var startDate, endDate *time.Time

	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if sd, err := time.Parse(time.RFC3339, startDateStr); err == nil {
			startDate = &sd
		}
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if ed, err := time.Parse(time.RFC3339, endDateStr); err == nil {
			endDate = &ed
		}
	}

	summary, err := h.auditLogService.GetAuditSummary(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to generate audit summary"})
		return
	}

	response := h.auditLogService.ConvertAuditLogSummaryToResponse(summary)
	c.JSON(http.StatusOK, response)
}

// CleanupOldLogs removes audit logs older than the retention period for LGPD compliance
func (h *AuditLogHandler) CleanupOldLogs(c *gin.Context) {
	retentionDays := 2555 // ~7 years for LGPD compliance
	if daysStr := c.Query("retention_days"); daysStr != "" {
		if days, err := strconv.Atoi(daysStr); err == nil && days > 0 {
			retentionDays = days
		}
	}

	err := h.auditLogService.CleanupOldLogs(retentionDays)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to cleanup old audit logs"})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{
		Message: "Audit logs cleanup completed",
	})
}
