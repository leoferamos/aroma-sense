package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/auth"
	"github.com/leoferamos/aroma-sense/internal/dto"
	"github.com/leoferamos/aroma-sense/internal/service"
)

type UserHandler struct {
	userService service.UserService
}

// NewUserHandler creates a new instance of UserHandler
func NewUserHandler(s service.UserService) *UserHandler {
	return &UserHandler{userService: s}
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
// @Failure      400  {object}  dto.ErrorResponse    "Invalid request (missing fields, invalid email) or email already registered"
// @Router       /users/register [post]
func (h *UserHandler) RegisterUser(c *gin.Context) {
	var input dto.CreateUserRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	input.Email = strings.ToLower(input.Email)

	if err := h.userService.RegisterUser(input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
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
// @Failure      400  {object}  dto.ErrorResponse  "Invalid request"
// @Router       /users/login [post]
func (h *UserHandler) LoginUser(c *gin.Context) {
	var input dto.LoginRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	input.Email = strings.ToLower(input.Email)

	accessToken, refreshToken, user, err := h.userService.Login(input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "invalid credentials"})
		return
	}

	// Set refresh token in HttpOnly cookie
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
// @Failure      401  {object}  dto.ErrorResponse   "Missing or invalid refresh token"
// @Router       /users/refresh [post]
func (h *UserHandler) RefreshToken(c *gin.Context) {
	// Read refresh token from HttpOnly cookie
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "missing refresh token"})
		return
	}

	// Validate refresh token and generate new access token + rotate refresh token
	accessToken, newRefreshToken, user, err := h.userService.RefreshAccessToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
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
		_ = h.userService.InvalidateRefreshToken(refreshToken)
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
// @Failure      401  {object}  dto.ErrorResponse      "Unauthorized"
// @Router       /users/me [get]
// @Security     BearerAuth
func (h *UserHandler) GetProfile(c *gin.Context) {
	rawUserID, exists := c.Get("userID")
	if !exists || rawUserID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}
	publicID := rawUserID.(string)
	user, err := h.userService.GetByPublicID(publicID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
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
// @Failure      400  {object}  dto.ErrorResponse       "Validation error"
// @Failure      401  {object}  dto.ErrorResponse       "Unauthorized"
// @Router       /users/me/profile [patch]
// @Security     BearerAuth
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	rawUserID, exists := c.Get("userID")
	if !exists || rawUserID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}
	publicID := rawUserID.(string)

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	user, err := h.userService.UpdateDisplayName(publicID, req.DisplayName)
	if err != nil {
		// simple validation mapping
		if strings.Contains(err.Error(), "too short") || strings.Contains(err.Error(), "too long") {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to update profile"})
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
// @Failure      401  {object}  dto.ErrorResponse      "Unauthorized"
// @Failure      500  {object}  dto.ErrorResponse      "Internal server error"
// @Router       /users/me/export [get]
// @Security     BearerAuth
func (h *UserHandler) ExportUserData(c *gin.Context) {
	rawUserID, exists := c.Get("userID")
	if !exists || rawUserID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}
	publicID := rawUserID.(string)

	data, err := h.userService.ExportUserData(publicID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to export user data"})
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
// @Failure      400      {object}  dto.ErrorResponse     "Invalid confirmation or active dependencies"
// @Failure      401      {object}  dto.ErrorResponse     "Unauthorized"
// @Router       /users/me/deletion [post]
// @Security     BearerAuth
func (h *UserHandler) RequestAccountDeletion(c *gin.Context) {
	rawUserID, exists := c.Get("userID")
	if !exists || rawUserID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}
	publicID := rawUserID.(string)

	var input dto.DeleteAccountRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	// Require explicit confirmation
	if input.Confirmation != "DELETE_MY_ACCOUNT" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid confirmation phrase"})
		return
	}

	if err := h.userService.RequestAccountDeletion(publicID); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "account deletion requested successfully - you have 7 days to change your mind"})
}

// ConfirmAccountDeletion confirms account deletion after cooling off period
//
// @Summary      Confirm account deletion
// @Description  Confirms account deletion after 7-day cooling off period has passed
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200      {object}  dto.MessageResponse   "Account deletion confirmed"
// @Failure      400      {object}  dto.ErrorResponse     "Cooling off period not expired or no deletion request"
// @Failure      401      {object}  dto.ErrorResponse     "Unauthorized"
// @Router       /users/me/deletion/confirm [post]
// @Security     BearerAuth
func (h *UserHandler) ConfirmAccountDeletion(c *gin.Context) {
	rawUserID, exists := c.Get("userID")
	if !exists || rawUserID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}
	publicID := rawUserID.(string)

	if err := h.userService.ConfirmAccountDeletion(publicID); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
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
// @Failure      400      {object}  dto.ErrorResponse     "No deletion request to cancel"
// @Failure      401      {object}  dto.ErrorResponse     "Unauthorized"
// @Router       /users/me/deletion/cancel [post]
// @Security     BearerAuth
func (h *UserHandler) CancelAccountDeletion(c *gin.Context) {
	rawUserID, exists := c.Get("userID")
	if !exists || rawUserID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}
	publicID := rawUserID.(string)

	if err := h.userService.CancelAccountDeletion(publicID); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
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
// @Failure      400  {object}  dto.ErrorResponse "Invalid request"
// @Failure      403  {object}  dto.ErrorResponse "No active deactivation or deadline expired"
// @Router       /users/me/contest [post]
// @Security     BearerAuth
func (h *UserHandler) RequestContestation(c *gin.Context) {
	rawUserID, exists := c.Get("userID")
	if !exists || rawUserID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}
	publicID := rawUserID.(string)

	var req dto.ContestationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.userService.RequestContestation(publicID, req.Reason); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "contestation request submitted successfully - our team will review it within 5 business days"})
}
