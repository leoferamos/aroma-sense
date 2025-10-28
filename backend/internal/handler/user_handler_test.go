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

func (m *MockUserService) InvalidateRefreshToken(refreshToken string) error {
	args := m.Called(refreshToken)
	return args.Error(0)
}

// ---- SETUP ROUTER ----
func setupUserRouter() (*gin.Engine, *MockUserService) {
	mockService := new(MockUserService)
	userHandler := handler.NewUserHandler(mockService)

	router := gin.Default()
	router.POST("/users/register", userHandler.RegisterUser)
	router.POST("/users/login", userHandler.LoginUser)

	return router, mockService
}

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
		router, mockService := setupUserRouter()
		payload := dto.CreateUserRequest{Email: "test@example.com", Password: "password123"}

		mockService.On("RegisterUser", payload).Return(nil)

		w := performRequest(t, router, "POST", "/users/register", payload)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "User registered successfully")
		mockService.AssertExpectations(t)
	})

	t.Run("Email Exists", func(t *testing.T) {
		router, mockService := setupUserRouter()
		payload := dto.CreateUserRequest{Email: "test@example.com", Password: "password123"}

		mockService.On("RegisterUser", payload).Return(errors.New("email already registered"))

		w := performRequest(t, router, "POST", "/users/register", payload)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "email already registered")
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid Payload", func(t *testing.T) {
		router, _ := setupUserRouter()
		w := performRequest(t, router, "POST", "/users/register", "invalid-payload")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestLoginUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		router, mockService := setupUserRouter()
		payload := dto.LoginRequest{Email: "test@example.com", Password: "password123"}
		expiresAt := time.Now().Add(7 * 24 * time.Hour)
		user := &model.User{
			PublicID:              "uuid",
			Email:                 "test@example.com",
			Role:                  "client",
			RefreshTokenExpiresAt: &expiresAt,
		}

		mockService.On("Login", payload).Return("mock-access", "mock-refresh", user, nil)

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
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid Credentials", func(t *testing.T) {
		router, mockService := setupUserRouter()
		payload := dto.LoginRequest{Email: "test@example.com", Password: "wrongpassword"}

	mockService.On("Login", payload).Return("", "", (*model.User)(nil), errors.New("invalid credentials"))

		w := performRequest(t, router, "POST", "/users/login", payload)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "invalid credentials")
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid Payload", func(t *testing.T) {
		router, _ := setupUserRouter()
		w := performRequest(t, router, "POST", "/users/login", "invalid-payload")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
