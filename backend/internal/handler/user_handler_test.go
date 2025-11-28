package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/dto"
	"github.com/leoferamos/aroma-sense/internal/handler"
	"github.com/leoferamos/aroma-sense/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ---- MOCK SERVICE ----
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) RegisterUser(input dto.CreateUserRequest) error {
	args := m.Called(input)
	return args.Error(0)
}

func (m *MockUserService) Login(input dto.LoginRequest) (string, string, *model.User, error) {
	args := m.Called(input)
	var user *model.User
	if args.Get(2) != nil {
		user = args.Get(2).(*model.User)
	}
	return args.String(0), args.String(1), user, args.Error(3)
}

func (m *MockUserService) RefreshAccessToken(refreshToken string) (string, string, *model.User, error) {
	args := m.Called(refreshToken)
	var user *model.User
	if args.Get(2) != nil {
		user = args.Get(2).(*model.User)
	}
	return args.String(0), args.String(1), user, args.Error(3)
}

func (m *MockUserService) Logout(refreshToken string) error {
	args := m.Called(refreshToken)
	return args.Error(0)
}

func (m *MockUserService) InvalidateRefreshToken(refreshToken string) error {
	args := m.Called(refreshToken)
	return args.Error(0)
}

func (m *MockUserService) GetByPublicID(publicID string) (*model.User, error) {
	args := m.Called(publicID)
	var user *model.User
	if args.Get(0) != nil {
		user = args.Get(0).(*model.User)
	}
	return user, args.Error(1)
}

func (m *MockUserService) UpdateDisplayName(publicID string, displayName string) (*model.User, error) {
	args := m.Called(publicID, displayName)
	var user *model.User
	if args.Get(0) != nil {
		user = args.Get(0).(*model.User)
	}
	return user, args.Error(1)
}

func (m *MockUserService) ListUsers(limit int, offset int, filters map[string]interface{}) ([]*model.User, int64, error) {
	args := m.Called(limit, offset, filters)
	var users []*model.User
	if args.Get(0) != nil {
		users = args.Get(0).([]*model.User)
	}
	return users, args.Get(1).(int64), args.Error(2)
}

func (m *MockUserService) GetUserByID(id uint) (*model.User, error) {
	args := m.Called(id)
	var user *model.User
	if args.Get(0) != nil {
		user = args.Get(0).(*model.User)
	}
	return user, args.Error(1)
}

func (m *MockUserService) UpdateUserRole(userID uint, newRole string, adminPublicID string) error {
	args := m.Called(userID, newRole, adminPublicID)
	return args.Error(0)
}

func (m *MockUserService) DeactivateUser(userID uint, adminPublicID string, reason string, notes string, suspensionUntil *time.Time) error {
	args := m.Called(userID, adminPublicID, reason, notes, suspensionUntil)
	return args.Error(0)
}

func (m *MockUserService) ExportUserData(publicID string) (*dto.UserExportResponse, error) {
	args := m.Called(publicID)
	var export *dto.UserExportResponse
	if args.Get(0) != nil {
		export = args.Get(0).(*dto.UserExportResponse)
	}
	return export, args.Error(1)
}

func (m *MockUserService) RequestAccountDeletion(publicID string) error {
	args := m.Called(publicID)
	return args.Error(0)
}

func (m *MockUserService) ConfirmAccountDeletion(publicID string) error {
	args := m.Called(publicID)
	return args.Error(0)
}

func (m *MockUserService) CancelAccountDeletion(publicID string) error {
	args := m.Called(publicID)
	return args.Error(0)
}

func (m *MockUserService) AnonymizeExpiredUser(publicID string) error {
	args := m.Called(publicID)
	return args.Error(0)
}

// ---- MOCKS PARA NOVA ASSINATURA ----
type MockAuthService struct{ mock.Mock }

func (m *MockAuthService) RegisterUser(input dto.CreateUserRequest) error { return nil }
func (m *MockAuthService) Login(input dto.LoginRequest) (string, string, *model.User, error) {
	return "", "", nil, nil
}
func (m *MockAuthService) RefreshAccessToken(refreshToken string) (string, string, *model.User, error) {
	return "", "", nil, nil
}
func (m *MockAuthService) Logout(refreshToken string) error                 { return nil }
func (m *MockAuthService) InvalidateRefreshToken(refreshToken string) error { return nil }

type MockUserProfileService struct{ mock.Mock }

func (m *MockUserProfileService) GetByPublicID(publicID string) (*model.User, error) { return nil, nil }
func (m *MockUserProfileService) UpdateDisplayName(publicID string, displayName string) (*model.User, error) {
	return nil, nil
}

type MockLgpdService struct{ mock.Mock }

func (m *MockLgpdService) ExportUserData(publicID string) (*dto.UserExportResponse, error) {
	return nil, nil
}
func (m *MockLgpdService) RequestAccountDeletion(publicID string) error             { return nil }
func (m *MockLgpdService) ConfirmAccountDeletion(publicID string) error             { return nil }
func (m *MockLgpdService) CancelAccountDeletion(publicID string) error              { return nil }
func (m *MockLgpdService) AnonymizeExpiredUser(publicID string) error               { return nil }
func (m *MockLgpdService) RequestContestation(publicID string, reason string) error { return nil }

// Adapte os métodos necessários para os testes, se necessário.

func setupUserRouter() (*gin.Engine, *MockAuthService, *MockUserProfileService, *MockLgpdService) {
	mockAuth := new(MockAuthService)
	mockProfile := new(MockUserProfileService)
	mockLgpd := new(MockLgpdService)
	userHandler := handler.NewUserHandler(mockAuth, mockProfile, mockLgpd)

	router := gin.Default()
	router.POST("/users/register", userHandler.RegisterUser)
	router.POST("/users/login", userHandler.LoginUser)
	router.GET("/users/me", userHandler.GetProfile)
	router.PATCH("/users/me/profile", userHandler.UpdateProfile)

	return router, mockAuth, mockProfile, mockLgpd
}
func TestGetProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockAuth := new(MockAuthService)
	mockProfile := new(MockUserProfileService)
	mockLgpd := new(MockLgpdService)
	user := &model.User{
		PublicID:    "uuid",
		Email:       "test@example.com",
		Role:        "client",
		DisplayName: ptr("Test User"),
		CreatedAt:   time.Now(),
	}
	mockProfile.On("GetByPublicID", "uuid").Return(user, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/users/me", nil)
	c.Set("userID", "uuid")
	handler := handler.NewUserHandler(mockAuth, mockProfile, mockLgpd)
	handler.GetProfile(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Test User")
	mockProfile.AssertExpectations(t)
}

func TestUpdateProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockAuth := new(MockAuthService)
	mockProfile := new(MockUserProfileService)
	mockLgpd := new(MockLgpdService)
	user := &model.User{
		PublicID:    "uuid",
		Email:       "test@example.com",
		Role:        "client",
		DisplayName: ptr("Updated Name"),
		CreatedAt:   time.Now(),
	}
	mockProfile.On("UpdateDisplayName", "uuid", "Updated Name").Return(user, nil)

	payload := dto.UpdateProfileRequest{DisplayName: "Updated Name"}
	body, _ := json.Marshal(payload)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PATCH", "/users/me/profile", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("userID", "uuid")
	handler := handler.NewUserHandler(mockAuth, mockProfile, mockLgpd)
	handler.UpdateProfile(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Updated Name")
	mockProfile.AssertExpectations(t)
}

func ptr(s string) *string { return &s }

func performRequest(t *testing.T, router *gin.Engine, method, url string, payload interface{}) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	var bodyBytes []byte
	if payload != nil {
		var err error
		bodyBytes, err = json.Marshal(payload)
		if err != nil {
			t.Fatalf("failed to marshal payload: %v", err)
		}
	}

	req, _ := http.NewRequest(method, url, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	return w
}

// ---- TESTS ----
func TestRegisterUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		router, mockAuth, _, _ := setupUserRouter()
		payload := dto.CreateUserRequest{Email: "test@example.com", Password: "StrongPass1"}

		mockAuth.On("RegisterUser", payload).Return(nil)

		w := performRequest(t, router, "POST", "/users/register", payload)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "User registered successfully")
		mockAuth.AssertExpectations(t)
	})

	t.Run("Email Exists", func(t *testing.T) {
		router, mockAuth, _, _ := setupUserRouter()
		payload := dto.CreateUserRequest{Email: "test@example.com", Password: "StrongPass1"}

		mockAuth.On("RegisterUser", payload).Return(errors.New("email already registered"))

		w := performRequest(t, router, "POST", "/users/register", payload)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "email already registered")
		mockAuth.AssertExpectations(t)
	})

	t.Run("Invalid Payload", func(t *testing.T) {
		router, _, _, _ := setupUserRouter()
		w := performRequest(t, router, "POST", "/users/register", "invalid-payload")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestLoginUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		router, mockAuth, _, _ := setupUserRouter()
		payload := dto.LoginRequest{Email: "test@example.com", Password: "StrongPass1"}
		expiresAt := time.Now().Add(7 * 24 * time.Hour)
		user := &model.User{
			PublicID:              "uuid",
			Email:                 "test@example.com",
			Role:                  "client",
			RefreshTokenExpiresAt: &expiresAt,
		}

		mockAuth.On("Login", payload).Return("mock-access", "mock-refresh", user, nil)

		w := performRequest(t, router, "POST", "/users/login", payload)

		assert.Equal(t, http.StatusOK, w.Code)

		// Check that refresh cookie is set
		cookies := w.Result().Cookies()
		var refreshCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == "refresh_token" {
				refreshCookie = cookie
				break
			}
		}
		assert.NotNil(t, refreshCookie, "Refresh cookie should be set")
		assert.Equal(t, "mock-refresh", refreshCookie.Value)

		// Check response body contains user data
		assert.Contains(t, w.Body.String(), "public_id")
		assert.Contains(t, w.Body.String(), "access_token")
		assert.Contains(t, w.Body.String(), "Login successful")
		mockAuth.AssertExpectations(t)
	})

	t.Run("Invalid Credentials", func(t *testing.T) {
		router, mockAuth, _, _ := setupUserRouter()
		payload := dto.LoginRequest{Email: "test@example.com", Password: "wrongpassword"}

		mockAuth.On("Login", payload).Return("", "", (*model.User)(nil), errors.New("invalid credentials"))

		w := performRequest(t, router, "POST", "/users/login", payload)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "invalid credentials")
		mockAuth.AssertExpectations(t)
	})

	t.Run("Invalid Payload", func(t *testing.T) {
		router, _, _, _ := setupUserRouter()
		w := performRequest(t, router, "POST", "/users/login", "invalid-payload")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
