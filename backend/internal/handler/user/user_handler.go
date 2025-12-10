package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/auth"
	"github.com/leoferamos/aroma-sense/internal/dto"
	handlererrors "github.com/leoferamos/aroma-sense/internal/handler/errors"
	"github.com/leoferamos/aroma-sense/internal/service"
)

type UserHandler struct {
	authService        service.AuthService
	userProfileService service.UserProfileService
	lgpdService        service.LgpdService
	chatService        service.ChatServiceInterface
}

// NewUserHandler creates a new instance of UserHandler
func NewUserHandler(auth service.AuthService, profile service.UserProfileService, lgpd service.LgpdService, chat service.ChatServiceInterface) *UserHandler {
	return &UserHandler{authService: auth, userProfileService: profile, lgpdService: lgpd, chatService: chat}
}

// UserProfile exposes the internal UserProfileService for middleware/router wiring
func (h *UserHandler) UserProfile() service.UserProfileService {
	return h.userProfileService
}

// RegisterUser handles user registration requests.
//
// @Summary      Register a new user
// @Description  Creates a new user account with the provided information.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body  dto.CreateUserRequest  true  "User registration data"
// @Success      201  {object}  dto.MessageResponse  "User registered successfully"
// @Failure      400  {object}  dto.ErrorResponse    "Error code: invalid_request or email_already_registered"
// @Router       /users/register [post]
func (h *UserHandler) RegisterUser(c *gin.Context) {
	var input dto.CreateUserRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request"})
		return
	}

	input.Email = strings.ToLower(input.Email)

	if err := h.authService.RegisterUser(input); err != nil {
		if status, code, ok := handlererrors.MapServiceError(err); ok {
			c.JSON(status, dto.ErrorResponse{Error: code})
			return
		}
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request"})
		return
	}

	c.JSON(http.StatusCreated, dto.MessageResponse{Message: "User registered successfully"})
}

// LoginUser handles user authentication requests.
//
// @Summary      Login
// @Description  Authenticates a user and returns a JWT token and user info.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body  dto.LoginRequest  true  "Login credentials"
// @Success      200  {object}  dto.LoginResponse   "Authentication successful"
// @Failure      400  {object}  dto.ErrorResponse  "Error code: invalid_request"
// @Router       /users/login [post]
func (h *UserHandler) LoginUser(c *gin.Context) {
	var input dto.LoginRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request"})
		return
	}

	input.Email = strings.ToLower(input.Email)

	accessToken, refreshToken, user, err := h.authService.Login(input)
	if err != nil {
		if status, code, ok := handlererrors.MapServiceError(err); ok {
			c.JSON(status, dto.ErrorResponse{Error: code})
			return
		}
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "invalid_credentials"})
		return
	}

	// Set refresh token in HttpOnly cookie
	if user.RefreshTokenExpiresAt == nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error"})
		return
	}
	auth.SetRefreshTokenCookie(c, refreshToken, *user.RefreshTokenExpiresAt)

	// Return access token in JSON response
	userResp := dto.UserResponse{
		PublicID:  user.PublicID,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}
	resp := dto.LoginResponse{
		Message:     "Login successful",
		AccessToken: accessToken,
		User:        userResp,
	}
	c.JSON(http.StatusOK, resp)
}

// RefreshToken generates a new access token and rotates the refresh token.
//
// @Summary      Refresh access token
// @Description  Uses the HttpOnly refresh_token cookie to issue a new short-lived access token.
// @Tags         auth
// @Produce      json
// @Success      200  {object}  dto.LoginResponse   "Token refreshed"
// @Failure      400  {object}  dto.ErrorResponse   "Error code: refresh_token_missing"
// @Failure      401  {object}  dto.ErrorResponse   "Error code: invalid_refresh_token or refresh_token_expired"
// @Router       /users/refresh [post]
func (h *UserHandler) RefreshToken(c *gin.Context) {
	// Read refresh token from HttpOnly cookie
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "refresh_token_missing"})
		return
	}

	// Validate refresh token and generate new access token + rotate refresh token
	accessToken, newRefreshToken, user, err := h.authService.RefreshAccessToken(refreshToken)
	if err != nil {
		if status, code, ok := handlererrors.MapServiceError(err); ok {
			c.JSON(status, dto.ErrorResponse{Error: code})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error"})
		return
	}

	// Set the rotated refresh token in HttpOnly cookie
	if user.RefreshTokenExpiresAt != nil {
		auth.SetRefreshTokenCookie(c, newRefreshToken, *user.RefreshTokenExpiresAt)
	}

	userResp := dto.UserResponse{
		PublicID:  user.PublicID,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}
	resp := dto.LoginResponse{
		Message:     "Token refreshed",
		AccessToken: accessToken,
		User:        userResp,
	}
	c.JSON(http.StatusOK, resp)
}

// LogoutUser handles user logout requests.
//
// @Summary      Logout
// @Description  Invalidates refresh token and clears refresh cookie.
// @Tags         auth
// @Success      200  {object}  dto.MessageResponse  "Logout successful"
// @Router       /users/logout [post]
func (h *UserHandler) LogoutUser(c *gin.Context) {
	if refreshToken, err := c.Cookie("refresh_token"); err == nil && refreshToken != "" {
		if err := h.authService.InvalidateRefreshToken(refreshToken); err != nil {
			if status, code, ok := handlererrors.MapServiceError(err); ok {
				c.JSON(status, dto.ErrorResponse{Error: code})
				return
			}
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error"})
			return
		}
	}
	auth.ClearRefreshTokenCookie(c)

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "Logout successful"})
}

// GetProfile returns the authenticated user's profile
//
// @Summary      Get my profile
// @Description  Returns the authenticated user's profile information.
// @Tags         users
// @Produce      json
// @Success      200  {object}  dto.ProfileResponse   "User profile"
// @Failure      401  {object}  dto.ErrorResponse      "Error code: unauthenticated"
// @Router       /users/me [get]
// @Security     BearerAuth
func (h *UserHandler) GetProfile(c *gin.Context) {
	publicID := c.GetString("userID")
	if publicID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthenticated"})
		return
	}
	user, err := h.userProfileService.GetByPublicID(publicID)
	if err != nil {
		if status, code, ok := handlererrors.MapServiceError(err); ok {
			c.JSON(status, dto.ErrorResponse{Error: code})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error"})
		return
	}
	resp := dto.ProfileResponse{
		PublicID:    user.PublicID,
		Email:       user.Email,
		Role:        user.Role,
		DisplayName: user.DisplayName,
		CreatedAt:   user.CreatedAt,
	}
	c.JSON(http.StatusOK, resp)
}

// UpdateProfile updates the authenticated user's display name
//
// @Summary      Update my profile
// @Description  Updates profile fields such as display_name.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        input  body  dto.UpdateProfileRequest  true  "Profile update"
// @Success      200  {object}  dto.ProfileResponse     "Updated profile"
// @Failure      400  {object}  dto.ErrorResponse       "Error code: invalid_request"
// @Failure      401  {object}  dto.ErrorResponse       "Error code: unauthenticated"
// @Router       /users/me/profile [patch]
// @Security     BearerAuth
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	publicID := c.GetString("userID")
	if publicID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthenticated"})
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request"})
		return
	}

	user, err := h.userProfileService.UpdateDisplayName(publicID, req.DisplayName)
	if err != nil {
		if status, code, ok := handlererrors.MapServiceError(err); ok {
			c.JSON(status, dto.ErrorResponse{Error: code})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error"})
		return
	}

	resp := dto.ProfileResponse{
		PublicID:    user.PublicID,
		Email:       user.Email,
		Role:        user.Role,
		DisplayName: user.DisplayName,
		CreatedAt:   user.CreatedAt,
	}
	c.JSON(http.StatusOK, resp)
}

// ExportUserData exports all user data for GDPR compliance
//
// @Summary      Export user data
// @Description  Download all personal data for GDPR portability right
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.UserExportResponse "User data exported"
// @Failure      401  {object}  dto.ErrorResponse      "Error code: unauthenticated"
// @Failure      500  {object}  dto.ErrorResponse      "Error code: internal_error"
// @Router       /users/me/export [get]
// @Security     BearerAuth
func (h *UserHandler) ExportUserData(c *gin.Context) {
	publicID := c.GetString("userID")
	if publicID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthenticated"})
		return
	}

	data, err := h.lgpdService.ExportUserData(publicID)
	if err != nil {
		if status, code, ok := handlererrors.MapServiceError(err); ok {
			c.JSON(status, dto.ErrorResponse{Error: code})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error"})
		return
	}

	c.JSON(http.StatusOK, data)
}

// RequestAccountDeletion initiates account deletion process with cooling off period
//
// @Summary      Request account deletion
// @Description  Initiates account deletion process with 7-day cooling off period (LGPD compliance)
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body   dto.DeleteAccountRequest true "Deletion confirmation"
// @Success      200      {object}  dto.MessageResponse   "Deletion request initiated successfully"
// @Failure      400      {object}  dto.ErrorResponse     "Error code: invalid_request or active_orders_block_deletion"
// @Failure      401      {object}  dto.ErrorResponse     "Error code: unauthenticated"
// @Router       /users/me/deletion [post]
// @Security     BearerAuth
func (h *UserHandler) RequestAccountDeletion(c *gin.Context) {
	publicID := c.GetString("userID")
	if publicID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthenticated"})
		return
	}

	var input dto.DeleteAccountRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request"})
		return
	}

	// Require explicit confirmation
	if input.Confirmation != "DELETE_MY_ACCOUNT" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request"})
		return
	}

	if err := h.lgpdService.RequestAccountDeletion(publicID); err != nil {
		if status, code, ok := handlererrors.MapServiceError(err); ok {
			c.JSON(status, dto.ErrorResponse{Error: code})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error"})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "account deletion requested successfully - you have 7 days to change your mind"})
}
func (h *UserHandler) ConfirmAccountDeletion(c *gin.Context) {
	publicID := c.GetString("userID")
	if publicID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthenticated"})
		return
	}

	if err := h.lgpdService.ConfirmAccountDeletion(publicID); err != nil {
		if status, code, ok := handlererrors.MapServiceError(err); ok {
			c.JSON(status, dto.ErrorResponse{Error: code})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error"})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "account deletion confirmed - your data will be retained for 2 years before permanent deletion"})
}

// CancelAccountDeletion cancels a pending account deletion request
//
// @Summary      Cancel account deletion
// @Description  Cancels a pending account deletion request during cooling off period
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200      {object}  dto.MessageResponse   "Account deletion cancelled"
// @Failure      400      {object}  dto.ErrorResponse     "Error code: deletion_not_requested"
// @Failure      401      {object}  dto.ErrorResponse     "Error code: unauthenticated"
// @Router       /users/me/deletion/cancel [post]
// @Security     BearerAuth
func (h *UserHandler) CancelAccountDeletion(c *gin.Context) {
	publicID := c.GetString("userID")
	if publicID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthenticated"})
		return
	}

	if err := h.lgpdService.CancelAccountDeletion(publicID); err != nil {
		if status, code, ok := handlererrors.MapServiceError(err); ok {
			c.JSON(status, dto.ErrorResponse{Error: code})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error"})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "account deletion cancelled successfully"})
}

// RequestContestation allows user to contest account deactivation (LGPD compliance)
//
// @Summary      Request account deactivation contestation
// @Description  User can contest their account deactivation within 7 days
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body   dto.ContestationRequest true "Contestation details"
// @Success      200  {object}  dto.MessageResponse "Contestation requested successfully"
// @Failure      400  {object}  dto.ErrorResponse "Error code: invalid_request"
// @Failure      403  {object}  dto.ErrorResponse "Error code: contestation_deadline_expired or account_not_deactivated"
// @Router       /users/me/contest [post]
// @Security     BearerAuth
func (h *UserHandler) RequestContestation(c *gin.Context) {
	publicID := c.GetString("userID")
	if publicID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthenticated"})
		return
	}

	var req dto.ContestationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request"})
		return
	}

	if err := h.lgpdService.RequestContestation(publicID, req.Reason); err != nil {
		if status, code, ok := handlererrors.MapServiceError(err); ok {
			c.JSON(status, dto.ErrorResponse{Error: code})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error"})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "contestation request submitted successfully - our team will review it within 5 business days"})
}

// ChangePassword allows authenticated users to change their password
// @Summary Change user password
// @Description Allows authenticated users to change their password by providing current and new password
// @Tags users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body dto.ChangePasswordRequest true "Password change request"
// @Success 200 {object} dto.MessageResponse "Password changed successfully"
// @Failure 400 {object} dto.ErrorResponse "Error code: invalid_request or current_password_incorrect"
// @Failure 401 {object} dto.ErrorResponse "Error code: unauthenticated"
// @Router /users/change-password [post]
// @Security     BearerAuth
func (h *UserHandler) ChangePassword(c *gin.Context) {
	publicID := c.GetString("userID")
	if publicID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthenticated"})
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request"})
		return
	}

	if err := h.userProfileService.ChangePassword(publicID, req.CurrentPassword, req.NewPassword); err != nil {
		if status, code, ok := handlererrors.MapServiceError(err); ok {
			c.JSON(status, dto.ErrorResponse{Error: code})
			return
		}
		if strings.Contains(err.Error(), "password must") {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request"})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error"})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "password changed successfully"})
}
