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

// ListAuditLogs
// @Summary      List audit logs
// @Description  Get paginated list of audit logs with filters (user, actor, action, resource, date)
// @Tags         audit-log
// @Accept       json
// @Produce      json
// @Param        user_id      query     int     false  "User ID"
// @Param        actor_id     query     int     false  "Actor ID"
// @Param        action       query     string  false  "Action type"
// @Param        resource     query     string  false  "Resource type"
// @Param        resource_id  query     string  false  "Resource ID"
// @Param        start_date   query     string  false  "Start date (RFC3339)"
// @Param        end_date     query     string  false  "End date (RFC3339)"
// @Param        limit        query     int     false  "Limit"
// @Param        offset       query     int     false  "Offset"
// @Success      200  {object}  dto.AuditLogListResponse
// @Failure      400  {object}  dto.ErrorResponse "Error code: invalid_request"
// @Failure      401  {object}  dto.ErrorResponse "Error code: unauthenticated"
// @Failure      500  {object}  dto.ErrorResponse "Error code: internal_error"
// @Router       /admin/audit-logs [get]
// @Security     BearerAuth
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
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error"})
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

// GetAuditLog
// @Summary      Get audit log
// @Description  Get a specific audit log entry (masked emails)
// @Tags         audit-log
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Audit Log ID"
// @Success      200  {object}  dto.AuditLogResponse
// @Failure      400  {object}  dto.ErrorResponse "Error code: invalid_request"
// @Failure      404  {object}  dto.ErrorResponse "Error code: not_found"
// @Router       /admin/audit-logs/{id} [get]
// @Security     BearerAuth
func (h *AuditLogHandler) GetAuditLog(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request"})
		return
	}

	auditLog, err := h.auditLogService.GetAuditLogByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "not_found"})
		return
	}

	response := h.auditLogService.ConvertAuditLogToResponse(auditLog)
	c.JSON(http.StatusOK, response)
}

// GetAuditLogDetailed
// @Summary      Get audit log (detailed)
// @Description  Get a specific audit log entry with full emails (admin/ops)
// @Tags         audit-log
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Audit Log ID"
// @Success      200  {object}  dto.AuditLogResponse
// @Failure      400  {object}  dto.ErrorResponse "Error code: invalid_request"
// @Failure      404  {object}  dto.ErrorResponse "Error code: not_found"
// @Router       /admin/audit-logs/{id}/detailed [get]
// @Security     BearerAuth
func (h *AuditLogHandler) GetAuditLogDetailed(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request"})
		return
	}

	auditLog, err := h.auditLogService.GetAuditLogByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "not_found"})
		return
	}

	response := h.auditLogService.ConvertAuditLogToResponseDetailed(auditLog)
	c.JSON(http.StatusOK, response)
}

// GetUserAuditLogs
// @Summary      Get user audit logs
// @Description  Get paginated audit logs for a specific user (masked emails)
// @Tags         audit-log
// @Accept       json
// @Produce      json
// @Param        id     path      int  true  "User ID"
// @Param        limit  query     int  false "Limit"
// @Param        offset query     int  false "Offset"
// @Success      200  {object}  dto.AuditLogListResponse
// @Failure      400  {object}  dto.ErrorResponse "Error code: invalid_request"
// @Failure      404  {object}  dto.ErrorResponse "Error code: not_found"
// @Router       /admin/users/{id}/audit-logs [get]
// @Security     BearerAuth
func (h *AuditLogHandler) GetUserAuditLogs(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request"})
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
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error"})
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

// GetAuditSummary
// @Summary      Get audit log summary
// @Description  Get summary/statistics of audit logs in a date range
// @Tags         audit-log
// @Accept       json
// @Produce      json
// @Param        start_date  query     string  false  "Start date (RFC3339)"
// @Param        end_date    query     string  false  "End date (RFC3339)"
// @Success      200  {object}  dto.AuditLogSummaryResponse
// @Failure      400  {object}  dto.ErrorResponse "Error code: invalid_request"
// @Failure      500  {object}  dto.ErrorResponse "Error code: internal_error"
// @Router       /admin/audit-logs/summary [get]
// @Security     BearerAuth
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
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error"})
		return
	}

	response := h.auditLogService.ConvertAuditLogSummaryToResponse(summary)
	c.JSON(http.StatusOK, response)
}

// CleanupOldLogs
// @Summary      Cleanup old audit logs
// @Description  Remove audit logs older than retention period (LGPD compliance)
// @Tags         audit-log
// @Accept       json
// @Produce      json
// @Param        retention_days  query     int  false  "Retention period in days"
// @Success      200  {object}  dto.MessageResponse
// @Failure      400  {object}  dto.ErrorResponse "Error code: invalid_request"
// @Failure      500  {object}  dto.ErrorResponse "Error code: internal_error"
// @Router       /admin/audit-logs/cleanup [post]
// @Security     BearerAuth
func (h *AuditLogHandler) CleanupOldLogs(c *gin.Context) {
	retentionDays := 2555 // ~7 years for LGPD compliance
	if daysStr := c.Query("retention_days"); daysStr != "" {
		if days, err := strconv.Atoi(daysStr); err == nil && days > 0 {
			retentionDays = days
		}
	}

	err := h.auditLogService.CleanupOldLogs(retentionDays)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error"})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{
		Message: "Audit logs cleanup completed",
	})
}
