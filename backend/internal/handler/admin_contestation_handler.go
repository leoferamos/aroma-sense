package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/dto"
	"github.com/leoferamos/aroma-sense/internal/service"
)

type AdminContestationHandler struct {
	service service.UserContestationService
}

func NewAdminContestationHandler(s service.UserContestationService) *AdminContestationHandler {
	return &AdminContestationHandler{service: s}
}

// ListPendingContestions returns all pending user contestations
//
// @Summary List pending user contestations
// @Description Returns all contestations with status 'pending' for admin review
// @Tags admin-contestation
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} map[string]interface{} "List of contestations and total count"
// @Failure 500 {object} dto.ErrorResponse
// @Router /admin/contestations [get]
// @Security BearerAuth
func (h *AdminContestationHandler) ListPendingContestions(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	contestations, total, err := h.service.ListPending(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}
	var dtos []dto.UserContestationResponse
	for _, cst := range contestations {
		dtos = append(dtos, dto.UserContestationResponseFromModel(&cst))
	}
	c.JSON(http.StatusOK, gin.H{"data": dtos, "total": total})
}

// ApproveContestation approves a pending contestation
//
// @Summary Approve a user contestation
// @Description Approves a pending contestation and adds optional review notes
// @Tags admin-contestation
// @Param id path int true "Contestation ID"
// @Param body body object false "Review notes (optional)"
// @Success 200 {object} dto.MessageResponse "Contestation approved"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /admin/contestations/{id}/approve [post]
// @Security BearerAuth
func (h *AdminContestationHandler) ApproveContestation(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}
	adminPublicID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}
	var req struct {
		Notes *string `json:"notes"`
	}
	_ = c.ShouldBindJSON(&req)
	if err := h.service.Approve(uint(id), adminPublicID.(string), req.Notes); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.MessageResponse{Message: "contestation approved"})
}

// RejectContestation rejects a pending contestation
//
// @Summary Reject a user contestation
// @Description Rejects a pending contestation and adds optional review notes
// @Tags admin-contestation
// @Param id path int true "Contestation ID"
// @Param body body object false "Review notes (optional)"
// @Success 200 {object} dto.MessageResponse "Contestation rejected"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /admin/contestations/{id}/reject [post]
// @Security BearerAuth
func (h *AdminContestationHandler) RejectContestation(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}
	adminPublicID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}
	var req struct {
		Notes *string `json:"notes"`
	}
	_ = c.ShouldBindJSON(&req)
	if err := h.service.Reject(uint(id), adminPublicID.(string), req.Notes); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.MessageResponse{Message: "contestation rejected"})
}
