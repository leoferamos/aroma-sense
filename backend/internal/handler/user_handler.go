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
