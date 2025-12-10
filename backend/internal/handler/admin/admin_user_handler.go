package admin

import (
	"net/http"
	"strconv"

	"github.com/leoferamos/aroma-sense/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/dto"
	handlererrors "github.com/leoferamos/aroma-sense/internal/handler/errors"
	"github.com/leoferamos/aroma-sense/internal/service"
)

type AdminUserHandler struct {
	userService service.AdminUserService
}

// NewAdminUserHandler creates a new instance of AdminUserHandler
func NewAdminUserHandler(s service.AdminUserService) *AdminUserHandler {
	return &AdminUserHandler{userService: s}
}

// AdminCreateAdmin allows a super admin to create a new admin user
func (h *AdminUserHandler) AdminCreateAdmin(c *gin.Context) {
	var input dto.AdminCreateUserRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request"})
		return
	}

	superID := c.GetString("userID")
	user, err := h.userService.CreateAdminUser(input.Email, input.Password, input.Name, superID)
	if err != nil {
		if status, code, ok := handlererrors.MapServiceError(err); ok {
			c.JSON(status, dto.ErrorResponse{Error: code})
			return
		}
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "admin created",
		"user": gin.H{
			"public_id": user.PublicID,
			"email":     user.Email,
			"role":      user.Role,
		},
	})
}

// AdminListUsers returns paginated list of users for admin
//
// @Summary      List users for admin
// @Description  Get paginated list of users with optional filters
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        limit   query     int                  false  "Number of users per page" default(10)
// @Param        offset  query     int                  false  "Offset for pagination" default(0)
// @Param        role    query     string               false  "Filter by role (admin, client)"
// @Param        status  query     string               false  "Filter by status (active, deactivated, deleted)"
// @Success      200     {object}  dto.UserListResponse "Users list"
// @Failure      400     {object}  dto.ErrorResponse    "Error code: invalid_request"
// @Failure      401     {object}  dto.ErrorResponse    "Error code: unauthenticated"
// @Failure      500     {object}  dto.ErrorResponse    "Error code: internal_error"
// @Router       /admin/users [get]
// @Security     BearerAuth
func (h *AdminUserHandler) AdminListUsers(c *gin.Context) {
	// Parse query parameters
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")
	role := c.Query("role")
	status := c.Query("status")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	// Build filters
	filters := make(map[string]interface{})
	if role != "" {
		filters["role"] = role
	}
	if status != "" {
		filters["status"] = status
	}

	// Get users
	users, total, err := h.userService.ListUsers(limit, offset, filters)
	if err != nil {
		if status, code, ok := handlererrors.MapServiceError(err); ok {
			c.JSON(status, dto.ErrorResponse{Error: code})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error"})
		return
	}

	// Convert to response DTO
	userResponses := make([]dto.AdminUserResponse, len(users))
	for i, user := range users {
		userResponses[i] = dto.AdminUserResponse{
			ID:                    user.ID,
			PublicID:              user.PublicID,
			MaskedEmail:           utils.MaskEmail(user.Email),
			Role:                  user.Role,
			DisplayName:           user.DisplayName,
			CreatedAt:             user.CreatedAt,
			LastLoginAt:           user.LastLoginAt,
			DeactivatedAt:         user.DeactivatedAt,
			DeactivatedBy:         user.DeactivatedBy,
			DeactivationReason:    user.DeactivationReason,
			DeactivationNotes:     user.DeactivationNotes,
			SuspensionUntil:       user.SuspensionUntil,
			ReactivationRequested: user.ReactivationRequested,
			ContestationDeadline:  user.ContestationDeadline,
		}
	}

	response := dto.UserListResponse{
		Users:  userResponses,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}

	c.JSON(http.StatusOK, response)
}

// AdminGetUser returns detailed user information for admin
//
// @Summary      Get user details for admin
// @Description  Get detailed information about a specific user
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        id   path      int                  true  "User ID"
// @Success      200  {object}  dto.AdminUserResponse "User details"
// @Failure      400  {object}  dto.ErrorResponse     "Error code: invalid_request"
// @Failure      401  {object}  dto.ErrorResponse     "Error code: unauthenticated"
// @Failure      404  {object}  dto.ErrorResponse     "Error code: not_found"
// @Failure      500  {object}  dto.ErrorResponse     "Error code: internal_error"
// @Router       /admin/users/{id} [get]
// @Security     BearerAuth
func (h *AdminUserHandler) AdminGetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request"})
		return
	}

	user, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		if status, code, ok := handlererrors.MapServiceError(err); ok {
			c.JSON(status, dto.ErrorResponse{Error: code})
			return
		}
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "not_found"})
		return
	}

	response := dto.AdminUserResponse{
		ID:                    user.ID,
		PublicID:              user.PublicID,
		MaskedEmail:           utils.MaskEmail(user.Email),
		Role:                  user.Role,
		DisplayName:           user.DisplayName,
		CreatedAt:             user.CreatedAt,
		LastLoginAt:           user.LastLoginAt,
		DeactivatedAt:         user.DeactivatedAt,
		DeactivatedBy:         user.DeactivatedBy,
		DeactivationReason:    user.DeactivationReason,
		DeactivationNotes:     user.DeactivationNotes,
		SuspensionUntil:       user.SuspensionUntil,
		ReactivationRequested: user.ReactivationRequested,
		ContestationDeadline:  user.ContestationDeadline,
	}

	c.JSON(http.StatusOK, response)
}

// AdminUpdateUserRole updates user role
//
// @Summary      Update user role
// @Description  Change user role (admin/client) with admin confirmation
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        id       path      int                  true  "User ID"
// @Param        role     body      dto.UpdateRoleRequest true "New role"
// @Success      200      {object}  dto.MessageResponse   "Role updated successfully"
// @Failure      400      {object}  dto.ErrorResponse     "Error code: invalid_request"
// @Failure      401      {object}  dto.ErrorResponse     "Error code: unauthenticated"
// @Failure      403      {object}  dto.ErrorResponse     "Error code: cannot_change_own_role"
// @Failure      404      {object}  dto.ErrorResponse     "Error code: not_found"
// @Router       /admin/users/{id}/role [patch]
// @Security     BearerAuth
func (h *AdminUserHandler) AdminUpdateUserRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request"})
		return
	}

	var input dto.UpdateRoleRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request"})
		return
	}

	// Get admin public ID from context
	adminPublicID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthenticated"})
		return
	}

	if err := h.userService.UpdateUserRole(uint(id), input.Role, adminPublicID.(string)); err != nil {
		if status, code, ok := handlererrors.MapServiceError(err); ok {
			c.JSON(status, dto.ErrorResponse{Error: code})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error"})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "User role updated successfully"})
}

// AdminDeactivateUser deactivates a user account with enhanced LGPD compliance
//
// @Summary      Deactivate user account
// @Description  Soft delete user account for GDPR compliance with detailed reason and notes
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        id   path      int                           true  "User ID"
// @Param        request body   dto.AdminDeactivateUserRequest true "Deactivation details"
// @Success      200  {object}  dto.MessageResponse           "User deactivated successfully"
// @Failure      400  {object}  dto.ErrorResponse             "Error code: invalid_request"
// @Failure      401  {object}  dto.ErrorResponse             "Error code: unauthenticated"
// @Failure      404  {object}  dto.ErrorResponse             "Error code: not_found"
// @Failure      500  {object}  dto.ErrorResponse             "Error code: internal_error"
// @Router       /admin/users/{id}/deactivate [post]
// @Security     BearerAuth
func (h *AdminUserHandler) AdminDeactivateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request"})
		return
	}

	// Get admin public ID from context
	adminPublicID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthenticated"})
		return
	}

	// Parse request body
	var req dto.AdminDeactivateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request"})
		return
	}

	if err := h.userService.DeactivateUser(uint(id), adminPublicID.(string), req.Reason, req.Notes, req.SuspensionUntil); err != nil {
		if status, code, ok := handlererrors.MapServiceError(err); ok {
			c.JSON(status, dto.ErrorResponse{Error: code})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error"})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "User deactivated successfully"})
}

// AdminReactivateUser reactivates a user account after contestation review
//
// @Summary      Reactivate user account
// @Description  Reactivate a previously deactivated user account with reason
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        id       path      int                        true  "User ID"
// @Param        request  body      dto.AdminReactivateUserRequest true "Reactivation details"
// @Success      200      {object}  dto.MessageResponse        "User reactivated successfully"
// @Failure      400      {object}  dto.ErrorResponse          "Error code: invalid_request"
// @Failure      401      {object}  dto.ErrorResponse          "Error code: unauthenticated"
// @Failure      404      {object}  dto.ErrorResponse          "Error code: not_found"
// @Failure      500      {object}  dto.ErrorResponse          "Error code: internal_error"
// @Router       /admin/users/{id}/reactivate [post]
// @Security     BearerAuth
func (h *AdminUserHandler) AdminReactivateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request"})
		return
	}

	// Get admin public ID from context
	adminPublicID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthenticated"})
		return
	}

	// Parse request body
	var req dto.AdminReactivateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request"})
		return
	}

	if err := h.userService.AdminReactivateUser(uint(id), adminPublicID.(string), req.Reason); err != nil {
		if status, code, ok := handlererrors.MapServiceError(err); ok {
			c.JSON(status, dto.ErrorResponse{Error: code})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error"})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "User reactivated successfully"})
}
