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
func (h *AdminContestationHandler) ApproveContestation(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}
	adminID, exists := c.Get("adminID")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}
	var req struct {
		Notes *string `json:"notes"`
	}
	_ = c.ShouldBindJSON(&req)
	if err := h.service.Approve(uint(id), adminID.(uint), req.Notes); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.MessageResponse{Message: "contestation approved"})
}

// RejectContestation rejects a pending contestation
func (h *AdminContestationHandler) RejectContestation(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}
	adminID, exists := c.Get("adminID")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}
	var req struct {
		Notes *string `json:"notes"`
	}
	_ = c.ShouldBindJSON(&req)
	if err := h.service.Reject(uint(id), adminID.(uint), req.Notes); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.MessageResponse{Message: "contestation rejected"})
}
