package admin_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/dto"
	"github.com/leoferamos/aroma-sense/internal/handler/admin"
	"github.com/leoferamos/aroma-sense/internal/model"
	serviceadmin "github.com/leoferamos/aroma-sense/internal/service/admin"
	"github.com/stretchr/testify/assert"
)

type mockAdminUserService struct {
	listUsersResult       []*model.User
	listUsersTotal        int64
	listUsersErr          error
	getUserByIDResult     *model.User
	getUserByIDErr        error
	updateUserRoleErr     error
	createAdminUserResult *model.User
	createAdminUserErr    error
	deactivateUserErr     error
	reactivateUserErr     error
}

func (m *mockAdminUserService) ListUsers(limit int, offset int, filters map[string]interface{}) ([]*model.User, int64, error) {
	return m.listUsersResult, m.listUsersTotal, m.listUsersErr
}

func (m *mockAdminUserService) GetUserByID(id uint) (*model.User, error) {
	return m.getUserByIDResult, m.getUserByIDErr
}

func (m *mockAdminUserService) UpdateUserRole(userID uint, newRole string, adminPublicID string) error {
	return m.updateUserRoleErr
}

func (m *mockAdminUserService) CreateAdminUser(email string, password string, displayName string, superAdminPublicID string) (*model.User, error) {
	return m.createAdminUserResult, m.createAdminUserErr
}

func (m *mockAdminUserService) DeactivateUser(userID uint, adminPublicID string, reason string, notes string, suspensionUntil *time.Time) error {
	return m.deactivateUserErr
}

func (m *mockAdminUserService) AdminReactivateUser(userID uint, adminPublicID string, reason string) error {
	return m.reactivateUserErr
}

// --- Test helpers ---
func createTestAdminUser() *model.User {
	now := time.Now()
	displayName := "Admin User"
	return &model.User{
		ID:          1,
		PublicID:    "user-123",
		Email:       "admin@example.com",
		Role:        "admin",
		DisplayName: &displayName,
		CreatedAt:   now,
		LastLoginAt: &now,
	}
}

func setupAdminUserRouter(svc serviceadmin.AdminUserService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Add middleware to simulate authentication
	r.Use(func(c *gin.Context) {
		c.Set("userID", "admin-123")
		c.Next()
	})

	handler := admin.NewAdminUserHandler(svc)
	r.POST("/admin/users/admin", handler.AdminCreateAdmin)
	r.GET("/admin/users", handler.AdminListUsers)
	r.GET("/admin/users/:id", handler.AdminGetUser)
	r.PATCH("/admin/users/:id/role", handler.AdminUpdateUserRole)
	r.POST("/admin/users/:id/deactivate", handler.AdminDeactivateUser)
	r.POST("/admin/users/:id/reactivate", handler.AdminReactivateUser)
	return r
}

// --- Tests ---
func TestAdminUserHandler_AdminCreateAdmin(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		user := createTestAdminUser()
		svc := &mockAdminUserService{
			createAdminUserResult: user,
		}
		r := setupAdminUserRouter(svc)

		reqBody := `{
			"email": "newadmin@example.com",
			"password": "securepassword123",
			"name": "New Admin"
		}`
		req, _ := http.NewRequest("POST", "/admin/users/admin", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "admin created", response["message"])
	})

	t.Run("invalid json", func(t *testing.T) {
		svc := &mockAdminUserService{}
		r := setupAdminUserRouter(svc)

		reqBody := `{"email": "invalid-email", "password": "pass"}`
		req, _ := http.NewRequest("POST", "/admin/users/admin", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid_request", response.Error)
	})

	t.Run("service error", func(t *testing.T) {
		svc := &mockAdminUserService{
			createAdminUserErr: errors.New("db error"),
		}
		r := setupAdminUserRouter(svc)

		reqBody := `{
			"email": "newadmin@example.com",
			"password": "securepassword123",
			"name": "New Admin"
		}`
		req, _ := http.NewRequest("POST", "/admin/users/admin", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid_request", response.Error)
	})
}

func TestAdminUserHandler_AdminListUsers(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		users := []*model.User{createTestAdminUser()}
		svc := &mockAdminUserService{
			listUsersResult: users,
			listUsersTotal:  1,
		}
		r := setupAdminUserRouter(svc)

		req, _ := http.NewRequest("GET", "/admin/users?limit=10&offset=0&role=admin&status=active", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.UserListResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Users, 1)
		assert.Equal(t, int64(1), response.Total)
		assert.Equal(t, 10, response.Limit)
		assert.Equal(t, 0, response.Offset)
	})

	t.Run("success with default params", func(t *testing.T) {
		users := []*model.User{createTestAdminUser()}
		svc := &mockAdminUserService{
			listUsersResult: users,
			listUsersTotal:  1,
		}
		r := setupAdminUserRouter(svc)

		req, _ := http.NewRequest("GET", "/admin/users", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.UserListResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 10, response.Limit) // default limit
		assert.Equal(t, 0, response.Offset) // default offset
	})

	t.Run("empty result", func(t *testing.T) {
		svc := &mockAdminUserService{
			listUsersResult: []*model.User{},
			listUsersTotal:  0,
		}
		r := setupAdminUserRouter(svc)

		req, _ := http.NewRequest("GET", "/admin/users", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.UserListResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Users, 0)
		assert.Equal(t, int64(0), response.Total)
	})

	t.Run("service error", func(t *testing.T) {
		svc := &mockAdminUserService{
			listUsersErr: errors.New("db error"),
		}
		r := setupAdminUserRouter(svc)

		req, _ := http.NewRequest("GET", "/admin/users", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "internal_error", response.Error)
	})
}

func TestAdminUserHandler_AdminGetUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		user := createTestAdminUser()
		svc := &mockAdminUserService{
			getUserByIDResult: user,
		}
		r := setupAdminUserRouter(svc)

		req, _ := http.NewRequest("GET", "/admin/users/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.AdminUserResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, uint(1), response.ID)
		assert.Equal(t, "user-123", response.PublicID)
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := &mockAdminUserService{}
		r := setupAdminUserRouter(svc)

		req, _ := http.NewRequest("GET", "/admin/users/invalid", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid_request", response.Error)
	})

	t.Run("user not found", func(t *testing.T) {
		svc := &mockAdminUserService{
			getUserByIDErr: errors.New("user not found"),
		}
		r := setupAdminUserRouter(svc)

		req, _ := http.NewRequest("GET", "/admin/users/999", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "not_found", response.Error)
	})

	t.Run("service error", func(t *testing.T) {
		svc := &mockAdminUserService{
			getUserByIDErr: errors.New("db error"),
		}
		r := setupAdminUserRouter(svc)

		req, _ := http.NewRequest("GET", "/admin/users/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "not_found", response.Error)
	})
}

func TestAdminUserHandler_AdminUpdateUserRole(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockAdminUserService{}
		r := setupAdminUserRouter(svc)

		reqBody := `{"role": "admin"}`
		req, _ := http.NewRequest("PATCH", "/admin/users/1/role", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.MessageResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "User role updated successfully", response.Message)
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := &mockAdminUserService{}
		r := setupAdminUserRouter(svc)

		reqBody := `{"role": "admin"}`
		req, _ := http.NewRequest("PATCH", "/admin/users/invalid/role", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid_request", response.Error)
	})

	t.Run("service error", func(t *testing.T) {
		svc := &mockAdminUserService{
			updateUserRoleErr: errors.New("db error"),
		}
		r := setupAdminUserRouter(svc)

		reqBody := `{"role": "admin"}`
		req, _ := http.NewRequest("PATCH", "/admin/users/1/role", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "internal_error", response.Error)
	})
}

func TestAdminUserHandler_AdminDeactivateUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockAdminUserService{}
		r := setupAdminUserRouter(svc)

		reqBody := `{
			"reason": "violation_of_terms",
			"notes": "User was abusive",
			"suspension_until": "2025-12-31T23:59:59Z"
		}`
		req, _ := http.NewRequest("POST", "/admin/users/1/deactivate", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.MessageResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "User deactivated successfully", response.Message)
	})

	t.Run("success minimal", func(t *testing.T) {
		svc := &mockAdminUserService{}
		r := setupAdminUserRouter(svc)

		reqBody := `{"reason": "violation_of_terms"}`
		req, _ := http.NewRequest("POST", "/admin/users/1/deactivate", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := &mockAdminUserService{}
		r := setupAdminUserRouter(svc)

		reqBody := `{"reason": "Violation of terms"}`
		req, _ := http.NewRequest("POST", "/admin/users/invalid/deactivate", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid_request", response.Error)
	})

	t.Run("invalid json", func(t *testing.T) {
		svc := &mockAdminUserService{}
		r := setupAdminUserRouter(svc)

		reqBody := `{"reason": "invalid_reason"}`
		req, _ := http.NewRequest("POST", "/admin/users/1/deactivate", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid_request", response.Error)
	})

	t.Run("service error", func(t *testing.T) {
		svc := &mockAdminUserService{
			deactivateUserErr: errors.New("db error"),
		}
		r := setupAdminUserRouter(svc)

		reqBody := `{"reason": "violation_of_terms"}`
		req, _ := http.NewRequest("POST", "/admin/users/1/deactivate", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "internal_error", response.Error)
	})
}

func TestAdminUserHandler_AdminReactivateUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockAdminUserService{}
		r := setupAdminUserRouter(svc)

		reqBody := `{"reason": "User appealed successfully"}`
		req, _ := http.NewRequest("POST", "/admin/users/1/reactivate", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.MessageResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "User reactivated successfully", response.Message)
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := &mockAdminUserService{}
		r := setupAdminUserRouter(svc)

		reqBody := `{"reason": "User appealed successfully"}`
		req, _ := http.NewRequest("POST", "/admin/users/invalid/reactivate", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid_request", response.Error)
	})

	t.Run("invalid json", func(t *testing.T) {
		svc := &mockAdminUserService{}
		r := setupAdminUserRouter(svc)

		reqBody := `{"reason": ""}`
		req, _ := http.NewRequest("POST", "/admin/users/1/reactivate", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid_request", response.Error)
	})

	t.Run("service error", func(t *testing.T) {
		svc := &mockAdminUserService{
			reactivateUserErr: errors.New("db error"),
		}
		r := setupAdminUserRouter(svc)

		reqBody := `{"reason": "User appealed successfully"}`
		req, _ := http.NewRequest("POST", "/admin/users/1/reactivate", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "internal_error", response.Error)
	})
}
