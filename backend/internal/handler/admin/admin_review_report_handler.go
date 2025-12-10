package admin

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/dto"
	handlererrors "github.com/leoferamos/aroma-sense/internal/handler"
	"github.com/leoferamos/aroma-sense/internal/service"
)

// AdminReviewReportHandler handles admin review report operations
type AdminReviewReportHandler struct {
	service service.ReviewReportService
}

func NewAdminReviewReportHandler(s service.ReviewReportService) *AdminReviewReportHandler {
	return &AdminReviewReportHandler{service: s}
}

// ListReports lists review reports filtered by status
//
// @Summary      List review reports
// @Description  List review reports filtered by status (pending/accepted/rejected)
// @Tags         admin-review-reports
// @Param        status  query    string  false  "Status filter"  Enums(pending,accepted,rejected)  default(pending)
// @Param        limit   query    int     false  "Limit"  default(20)
// @Param        offset  query    int     false  "Offset" default(0)
// @Success      200  {object}  dto.ReviewReportAdminResponse
// @Failure      400  {object}  dto.ErrorResponse "Error code: invalid_status"
// @Failure      401  {object}  dto.ErrorResponse "Error code: unauthenticated"
// @Failure      403  {object}  dto.ErrorResponse "Error code: unauthorized"
// @Failure      500  {object}  dto.ErrorResponse "Error code: internal_error"
// @Router       /admin/review-reports [get]
// @Security     BearerAuth
func (h *AdminReviewReportHandler) ListReports(c *gin.Context) {
	status := c.DefaultQuery("status", "pending")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	reports, total, err := h.service.List(c.Request.Context(), status, limit, offset)
	if err != nil {
		if statusCode, code, ok := handlererrors.MapServiceError(err); ok {
			c.JSON(statusCode, dto.ErrorResponse{Error: code})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error"})
		return
	}

	items := make([]dto.ReviewReportAdminItem, 0, len(reports))
	for i := range reports {
		items = append(items, dto.ReviewReportAdminItemFromModel(&reports[i]))
	}

	c.JSON(http.StatusOK, dto.ReviewReportAdminResponse{
		Items:  items,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

// ResolveReport resolves a review report (accept or reject)
//
// @Summary      Resolve a review report
// @Description  Accept or reject a review report. Accepting hides the review. Optionally deactivate the user later.
// @Tags         admin-review-reports
// @Param        id      path     string                           true  "Report ID"
// @Param        body    body     dto.ReviewReportResolveRequest   true  "Resolve payload"
// @Success      200  {object}  dto.MessageResponse
// @Failure      400  {object}  dto.ErrorResponse "Error code: invalid_action or invalid_request"
// @Failure      401  {object}  dto.ErrorResponse "Error code: unauthenticated"
// @Failure      403  {object}  dto.ErrorResponse "Error code: unauthorized"
// @Failure      404  {object}  dto.ErrorResponse "Error code: report_not_found or review_not_found"
// @Failure      409  {object}  dto.ErrorResponse "Error code: report_already_resolved"
// @Failure      500  {object}  dto.ErrorResponse "Error code: internal_error"
// @Router       /admin/review-reports/{id}/resolve [post]
// @Security     BearerAuth
func (h *AdminReviewReportHandler) ResolveReport(c *gin.Context) {
	reportID := c.Param("id")

	adminPublicID, exists := c.Get("userID")
	if !exists || adminPublicID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthenticated"})
		return
	}

	var req dto.ReviewReportResolveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request"})
		return
	}

	var suspensionUntil *time.Time
	if req.SuspensionUntil != nil && strings.TrimSpace(*req.SuspensionUntil) != "" {
		parsed, err := time.Parse(time.RFC3339, *req.SuspensionUntil)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request"})
			return
		}
		suspensionUntil = &parsed
	}

	if err := h.service.Resolve(c.Request.Context(), reportID, req.Action, req.DeactivateUser, adminPublicID.(string), suspensionUntil, req.Notes); err != nil {
		if statusCode, code, ok := handlererrors.MapServiceError(err); ok {
			c.JSON(statusCode, dto.ErrorResponse{Error: code})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error"})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "report resolved"})
}
