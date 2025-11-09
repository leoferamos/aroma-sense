package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/dto"
	"github.com/leoferamos/aroma-sense/internal/service"
)

type PasswordResetHandler struct {
	resetService service.PasswordResetService
}

// NewPasswordResetHandler creates a new instance of PasswordResetHandler
func NewPasswordResetHandler(s service.PasswordResetService) *PasswordResetHandler {
	return &PasswordResetHandler{resetService: s}
}

// RequestReset handles password reset requests.
// Always returns success to prevent email enumeration attacks.
//
// @Summary      Request password reset
// @Description  Sends a 6-digit code to the user's email if it exists. Always returns success to avoid revealing whether the email is registered.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body  dto.ResetPasswordRequestRequest  true  "Email address"
// @Success      200  {object}  dto.MessageResponse  "If the email exists, a reset code has been sent"
// @Failure      400  {object}  dto.ErrorResponse    "Invalid request (missing or malformed email)"
// @Failure      500  {object}  dto.ErrorResponse    "Internal server error"
// @Router       /users/reset/request [post]
func (h *PasswordResetHandler) RequestReset(c *gin.Context) {
	var input dto.ResetPasswordRequestRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid email format"})
		return
	}

	// Normalize email
	input.Email = strings.TrimSpace(strings.ToLower(input.Email))

	// Call service
	if err := h.resetService.RequestReset(input.Email); err != nil {
		// Log error internally but don't expose details to user
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to process request"})
		return
	}

	// Generic success message
	c.JSON(http.StatusOK, dto.MessageResponse{
		Message: "If the email exists, a reset code has been sent to your inbox",
	})
}

// ConfirmReset handles password reset confirmation with code.
// Validates the code and resets the password if valid.
//
// @Summary      Confirm password reset
// @Description  Validates the reset code and updates the user's password. Returns generic errors to avoid revealing whether email/code exists.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body  dto.ResetPasswordConfirmRequest  true  "Reset confirmation data"
// @Success      200  {object}  dto.MessageResponse  "Password reset successfully"
// @Failure      400  {object}  dto.ErrorResponse    "Invalid request, invalid/expired code, or weak password"
// @Failure      500  {object}  dto.ErrorResponse    "Internal server error"
// @Router       /users/reset/confirm [post]
func (h *PasswordResetHandler) ConfirmReset(c *gin.Context) {
	var input dto.ResetPasswordConfirmRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid request format"})
		return
	}

	// Normalize email and code
	input.Email = strings.TrimSpace(strings.ToLower(input.Email))
	input.Code = strings.TrimSpace(input.Code)

	// Validate code format
	if len(input.Code) != 6 {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid or expired reset code"})
		return
	}

	// Call service to reset password
	if err := h.resetService.ConfirmReset(input.Email, input.Code, input.NewPassword); err != nil {
		// Return generic error to avoid revealing details
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{
		Message: "Password reset successfully",
	})
}
