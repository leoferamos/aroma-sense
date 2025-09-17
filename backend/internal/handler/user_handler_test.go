package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

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

func (m *MockUserService) Login(input dto.LoginRequest) (string, *model.User, error) {
	args := m.Called(input)
	var user *model.User
	if args.Get(1) != nil {
		user = args.Get(1).(*model.User)
	}
	return args.String(0), user, args.Error(2)
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
		user := &model.User{PublicID: "uuid", Email: "test@example.com", Role: "client"}

		mockService.On("Login", payload).Return("mock-token", user, nil)

		w := performRequest(t, router, "POST", "/users/login", payload)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "mock-token")
		assert.Contains(t, w.Body.String(), "public_id")
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid Credentials", func(t *testing.T) {
		router, mockService := setupUserRouter()
		payload := dto.LoginRequest{Email: "test@example.com", Password: "wrongpassword"}

		mockService.On("Login", payload).Return("", (*model.User)(nil), errors.New("invalid credentials"))

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
