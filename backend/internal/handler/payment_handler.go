package handler

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/dto"
	"github.com/leoferamos/aroma-sense/internal/service"
)

// PaymentHandler handles payment-related endpoints.
type PaymentHandler struct {
	paymentService service.PaymentService
}

func NewPaymentHandler(paymentService service.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService}
}

// CreateIntent creates a PaymentIntent and returns the client secret.
func (h *PaymentHandler) CreateIntent(c *gin.Context) {
	var req dto.CreatePaymentIntentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request"})
		return
	}

	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthenticated"})
		return
	}

	res, err := h.paymentService.CreateIntent(c.Request.Context(), userID, &req)
	if err != nil {
		if status, code, ok := mapServiceError(err); ok {
			c.JSON(status, dto.ErrorResponse{Error: code})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"payment_intent_id": res.ID,
		"client_secret":     res.ClientSecret,
	})
}

// HandleWebhook will validate and process payment webhooks (stub for now).
func (h *PaymentHandler) HandleWebhook(c *gin.Context) {
	signature := c.GetHeader("Stripe-Signature")
	if signature == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "missing_signature"})
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error"})
		return
	}

	if _, err := h.paymentService.HandleWebhook(c.Request.Context(), body, signature); err != nil {
		if status, code, ok := mapServiceError(err); ok {
			c.JSON(status, dto.ErrorResponse{Error: code})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error"})
		return
	}

	c.Status(http.StatusOK)
}
