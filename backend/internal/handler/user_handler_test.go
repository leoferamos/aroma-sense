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
)

// ---- MOCK SERVICE ----
type mockUserService struct {
	RegisterUserFunc func(input dto.CreateUserRequest) error
	LoginFunc        func(input dto.LoginRequest) (string, *model.User, error)
}

func (m *mockUserService) RegisterUser(input dto.CreateUserRequest) error {
	return m.RegisterUserFunc(input)
}
func (m *mockUserService) Login(input dto.LoginRequest) (string, *model.User, error) {
	return m.LoginFunc(input)
}

// ---- ROUTE HANDLER MAP ----
type routeHandler func(*gin.Context)
type routeMap map[string]routeHandler

func newRouteMap(h *handler.UserHandler) routeMap {
	return routeMap{
		"/users/register": h.RegisterUser,
		"/users/login":    h.LoginUser,
	}
}

// ---- HELPER ----
func performRequest(t *testing.T, routes routeMap, method, url string, payload interface{}) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	var bodyBytes []byte
	var err error
	switch p := payload.(type) {
	case string:
		bodyBytes = []byte(p)
	default:
		bodyBytes, err = json.Marshal(payload)
		if err != nil {
			t.Fatalf("failed to marshal payload: %v", err)
		}
	}

	c.Request = httptest.NewRequest(method, url, bytes.NewBuffer(bodyBytes))
	c.Request.Header.Set("Content-Type", "application/json")

	handlerFunc, ok := routes[url]
	if !ok {
		t.Fatalf("route not found: %s", url)
	}

	handlerFunc(c)
	return w
}

// ---- TESTS ----
func TestRegisterUser_Success(t *testing.T) {
	mockSvc := &mockUserService{RegisterUserFunc: func(input dto.CreateUserRequest) error { return nil }}
	h := handler.NewUserHandler(mockSvc)
	routes := newRouteMap(h)

	payload := dto.CreateUserRequest{Email: "test@example.com", Password: "12345678"}
	w := performRequest(t, routes, "POST", "/users/register", payload)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "User registered successfully")
}

func TestRegisterUser_EmailExists(t *testing.T) {
	mockSvc := &mockUserService{RegisterUserFunc: func(input dto.CreateUserRequest) error { return errors.New("email already registered") }}
	h := handler.NewUserHandler(mockSvc)
	routes := newRouteMap(h)

	payload := dto.CreateUserRequest{Email: "test@example.com", Password: "12345678"}
	w := performRequest(t, routes, "POST", "/users/register", payload)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "email already registered")
}

func TestLoginUser_Success(t *testing.T) {
	mockSvc := &mockUserService{
		LoginFunc: func(input dto.LoginRequest) (string, *model.User, error) {
			return "mocktoken", &model.User{PublicID: "uuid", Email: "test@example.com", Role: "client"}, nil
		},
	}
	h := handler.NewUserHandler(mockSvc)
	routes := newRouteMap(h)

	payload := dto.LoginRequest{Email: "test@example.com", Password: "12345678"}
	w := performRequest(t, routes, "POST", "/users/login", payload)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "mocktoken")
	assert.Contains(t, w.Body.String(), "public_id")
}

func TestLoginUser_InvalidCredentials(t *testing.T) {
	mockSvc := &mockUserService{LoginFunc: func(input dto.LoginRequest) (string, *model.User, error) {
		return "", nil, errors.New("invalid credentials")
	}}
	h := handler.NewUserHandler(mockSvc)
	routes := newRouteMap(h)

	payload := dto.LoginRequest{Email: "test@example.com", Password: "wrongpass"}
	w := performRequest(t, routes, "POST", "/users/login", payload)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid credentials")
}

func TestRegisterUser_InvalidPayload(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	mockSvc := &mockUserService{RegisterUserFunc: func(input dto.CreateUserRequest) error { return nil }}
	h := handler.NewUserHandler(mockSvc)

	// Invalid Payload
	body := []byte(`{"password": "12345678"}`)
	c.Request = httptest.NewRequest("POST", "/users/register", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	h.RegisterUser(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

func TestLoginUser_InvalidPayload(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	mockSvc := &mockUserService{LoginFunc: func(input dto.LoginRequest) (string, *model.User, error) { return "", nil, nil }}
	h := handler.NewUserHandler(mockSvc)

	// Invalid Payload
	body := []byte(`{"password": "12345678"}`)
	c.Request = httptest.NewRequest("POST", "/users/login", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	h.LoginUser(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}
