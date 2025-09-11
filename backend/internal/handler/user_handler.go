// internal/handler/user_handler.go
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
	token, user, err := h.userService.Login(input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "invalid credentials"})
		return
	}

	userResp := dto.UserResponse{
		PublicID:  user.PublicID,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	c.JSON(http.StatusOK, dto.LoginResponse{Token: token, User: userResp})
}
