package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/dto"
	"github.com/leoferamos/aroma-sense/internal/rate"
	"github.com/leoferamos/aroma-sense/internal/service"
)

type PasswordResetHandler struct {
	resetService service.PasswordResetService
	rateLimiter  rate.RateLimiter
}

// NewPasswordResetHandler creates a new instance of PasswordResetHandler
func NewPasswordResetHandler(s service.PasswordResetService, limiter rate.RateLimiter) *PasswordResetHandler {
	return &PasswordResetHandler{
		resetService: s,
		rateLimiter:  limiter,
	}
}

// RequestReset handles password reset requests.
// Always returns success to prevent email enumeration attacks.
// Rate limited to 3 requests per email per hour.
//
// @Summary      Request password reset
// @Description  Sends a 6-digit code to the user's email if it exists. Always returns success to avoid revealing whether the email is registered.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body  dto.ResetPasswordRequestRequest  true  "Email address"
// @Success      200  {object}  dto.MessageResponse  "If the email exists, a reset code has been sent"
// @Failure      400  {object}  dto.ErrorResponse    "Error code: invalid_request"
// @Failure      429  {object}  dto.ErrorResponse    "Error code: rate_limited"
// @Failure      500  {object}  dto.ErrorResponse    "Error code: internal_error"
// @Router       /users/reset/request [post]
func (h *PasswordResetHandler) RequestReset(c *gin.Context) {
	var input dto.ResetPasswordRequestRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request"})
		return
	}

	// Normalize email
	input.Email = strings.TrimSpace(strings.ToLower(input.Email))

	// Apply rate limiting: 3 requests per email per hour
	ctx := context.Background()
	bucket := fmt.Sprintf("reset_request:%s", input.Email)
	allowed, remaining, resetAt, err := h.rateLimiter.Allow(ctx, bucket, 3, time.Hour)

	// Add rate limit headers
	c.Header("X-RateLimit-Limit", "3")
	c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
	c.Header("X-RateLimit-Reset", resetAt.Format(time.RFC3339))

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error"})
		return
	}

	if !allowed {
		if resetAt.After(time.Now()) {
			retryAfter := int(time.Until(resetAt).Seconds())
			if retryAfter < 0 {
				retryAfter = 0
			}
			c.Header("Retry-After", strconv.Itoa(retryAfter))
		}
		c.JSON(http.StatusTooManyRequests, dto.ErrorResponse{Error: "rate_limited"})
		return
	}

	// Call service
	if err := h.resetService.RequestReset(input.Email); err != nil {
		// Log error internally but don't expose details to user
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error"})
		return
	}

	// Generic success message
	c.JSON(http.StatusOK, dto.MessageResponse{
		Message: "If the email exists, a reset code has been sent to your inbox",
	})
}

// ConfirmReset handles password reset confirmation with code.
// Validates the code and resets the password if valid.
// Rate limited to 10 attempts per IP per hour.
//
// @Summary      Confirm password reset
// @Description  Validates the reset code and updates the user's password. Returns generic errors to avoid revealing whether email/code exists.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body  dto.ResetPasswordConfirmRequest  true  "Reset confirmation data"
// @Success      200  {object}  dto.MessageResponse  "Password reset successfully"
// @Failure      400  {object}  dto.ErrorResponse    "Error code: invalid_request or reset_code_invalid"
// @Failure      429  {object}  dto.ErrorResponse    "Error code: rate_limited"
// @Failure      500  {object}  dto.ErrorResponse    "Error code: internal_error"
// @Router       /users/reset/confirm [post]
func (h *PasswordResetHandler) ConfirmReset(c *gin.Context) {
	var input dto.ResetPasswordConfirmRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request"})
		return
	}

	// Normalize email and code
	input.Email = strings.TrimSpace(strings.ToLower(input.Email))
	input.Code = strings.TrimSpace(input.Code)

	// Validate code format
	if len(input.Code) != 6 {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request"})
		return
	}

	// Apply rate limiting: 10 attempts per IP per hour
	ctx := context.Background()
	clientIP := c.ClientIP()
	bucket := fmt.Sprintf("reset_confirm:%s", clientIP)
	allowed, remaining, resetAt, err := h.rateLimiter.Allow(ctx, bucket, 10, time.Hour)

	// Add rate limit headers
	c.Header("X-RateLimit-Limit", "10")
	c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
	c.Header("X-RateLimit-Reset", resetAt.Format(time.RFC3339))

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error"})
		return
	}

	if !allowed {
		if resetAt.After(time.Now()) {
			retryAfter := int(time.Until(resetAt).Seconds())
			if retryAfter < 0 {
				retryAfter = 0
			}
			c.Header("Retry-After", strconv.Itoa(retryAfter))
		}
		c.JSON(http.StatusTooManyRequests, dto.ErrorResponse{Error: "rate_limited"})
		return
	}

	// Call service to reset password
	if err := h.resetService.ConfirmReset(input.Email, input.Code, input.NewPassword); err != nil {
		// Return generic error to avoid revealing details
		if status, code, ok := mapServiceError(err); ok {
			c.JSON(status, dto.ErrorResponse{Error: code})
			return
		}
		if strings.Contains(err.Error(), "password must") {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request"})
			return
		}
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request"})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{
		Message: "Password reset successfully",
	})
}
