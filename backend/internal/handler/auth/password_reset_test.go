package auth_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/dto"
	"github.com/leoferamos/aroma-sense/internal/handler/auth"
	"github.com/leoferamos/aroma-sense/internal/rate"
	authservice "github.com/leoferamos/aroma-sense/internal/service/auth"
	"github.com/stretchr/testify/assert"
)

type mockPasswordResetService struct {
	requestResetErr error
	confirmResetErr error
}

func (m *mockPasswordResetService) RequestReset(email string) error {
	return m.requestResetErr
}

func (m *mockPasswordResetService) ConfirmReset(email string, code string, newPassword string) error {
	return m.confirmResetErr
}

type mockRateLimiter struct {
	allowResult    bool
	allowRemaining int
	allowResetAt   time.Time
	allowErr       error
}

func (m *mockRateLimiter) Allow(ctx context.Context, bucket string, limit int, window time.Duration) (bool, int, time.Time, error) {
	return m.allowResult, m.allowRemaining, m.allowResetAt, m.allowErr
}

func setupPasswordResetRouter(svc authservice.PasswordResetService, limiter rate.RateLimiter) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	handler := auth.NewPasswordResetHandler(svc, limiter)
	r.POST("/users/reset/request", handler.RequestReset)
	r.POST("/users/reset/confirm", handler.ConfirmReset)
	return r
}

// --- Tests ---
func TestPasswordResetHandler_RequestReset(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockPasswordResetService{}
		limiter := &mockRateLimiter{
			allowResult:    true,
			allowRemaining: 2,
			allowResetAt:   time.Now().Add(time.Hour),
		}
		r := setupPasswordResetRouter(svc, limiter)

		reqBody := `{"email": "user@example.com"}`
		req, _ := http.NewRequest("POST", "/users/reset/request", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.MessageResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "If the email exists, a reset code has been sent to your inbox", response.Message)

		// Check rate limit headers
		assert.Equal(t, "3", w.Header().Get("X-RateLimit-Limit"))
		assert.Equal(t, "2", w.Header().Get("X-RateLimit-Remaining"))
	})

	t.Run("invalid json", func(t *testing.T) {
		svc := &mockPasswordResetService{}
		limiter := &mockRateLimiter{
			allowResult:    true,
			allowRemaining: 2,
			allowResetAt:   time.Now().Add(time.Hour),
		}
		r := setupPasswordResetRouter(svc, limiter)

		reqBody := `{"email": ""}`
		req, _ := http.NewRequest("POST", "/users/reset/request", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid_request", response.Error)
	})

	t.Run("rate limited", func(t *testing.T) {
		svc := &mockPasswordResetService{}
		resetAt := time.Now().Add(30 * time.Minute)
		limiter := &mockRateLimiter{
			allowResult:    false,
			allowRemaining: 0,
			allowResetAt:   resetAt,
		}
		r := setupPasswordResetRouter(svc, limiter)

		reqBody := `{"email": "user@example.com"}`
		req, _ := http.NewRequest("POST", "/users/reset/request", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusTooManyRequests, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "rate_limited", response.Error)

		// Check rate limit headers
		assert.Equal(t, "3", w.Header().Get("X-RateLimit-Limit"))
		assert.Equal(t, "0", w.Header().Get("X-RateLimit-Remaining"))
		assert.NotEmpty(t, w.Header().Get("Retry-After"))
	})

	t.Run("rate limiter error", func(t *testing.T) {
		svc := &mockPasswordResetService{}
		limiter := &mockRateLimiter{
			allowErr: errors.New("redis error"),
		}
		r := setupPasswordResetRouter(svc, limiter)

		reqBody := `{"email": "user@example.com"}`
		req, _ := http.NewRequest("POST", "/users/reset/request", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "internal_error", response.Error)
	})

	t.Run("service error", func(t *testing.T) {
		svc := &mockPasswordResetService{
			requestResetErr: errors.New("email service error"),
		}
		limiter := &mockRateLimiter{
			allowResult:    true,
			allowRemaining: 2,
			allowResetAt:   time.Now().Add(time.Hour),
		}
		r := setupPasswordResetRouter(svc, limiter)

		reqBody := `{"email": "user@example.com"}`
		req, _ := http.NewRequest("POST", "/users/reset/request", strings.NewReader(reqBody))
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

func TestPasswordResetHandler_ConfirmReset(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockPasswordResetService{}
		limiter := &mockRateLimiter{
			allowResult:    true,
			allowRemaining: 9,
			allowResetAt:   time.Now().Add(time.Hour),
		}
		r := setupPasswordResetRouter(svc, limiter)

		reqBody := `{
			"email": "user@example.com",
			"code": "123456",
			"new_password": "newpassword123"
		}`
		req, _ := http.NewRequest("POST", "/users/reset/confirm", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.MessageResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Password reset successfully", response.Message)

		// Check rate limit headers
		assert.Equal(t, "10", w.Header().Get("X-RateLimit-Limit"))
		assert.Equal(t, "9", w.Header().Get("X-RateLimit-Remaining"))
	})

	t.Run("invalid json", func(t *testing.T) {
		svc := &mockPasswordResetService{}
		limiter := &mockRateLimiter{
			allowResult:    true,
			allowRemaining: 9,
			allowResetAt:   time.Now().Add(time.Hour),
		}
		r := setupPasswordResetRouter(svc, limiter)

		reqBody := `{"email": "", "code": "123456", "new_password": "pass"}`
		req, _ := http.NewRequest("POST", "/users/reset/confirm", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid_request", response.Error)
	})

	t.Run("invalid code length", func(t *testing.T) {
		svc := &mockPasswordResetService{}
		limiter := &mockRateLimiter{
			allowResult:    true,
			allowRemaining: 9,
			allowResetAt:   time.Now().Add(time.Hour),
		}
		r := setupPasswordResetRouter(svc, limiter)

		reqBody := `{
			"email": "user@example.com",
			"code": "12345",
			"new_password": "newpassword123"
		}`
		req, _ := http.NewRequest("POST", "/users/reset/confirm", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid_request", response.Error)
	})

	t.Run("rate limited", func(t *testing.T) {
		svc := &mockPasswordResetService{}
		resetAt := time.Now().Add(15 * time.Minute)
		limiter := &mockRateLimiter{
			allowResult:    false,
			allowRemaining: 0,
			allowResetAt:   resetAt,
		}
		r := setupPasswordResetRouter(svc, limiter)

		reqBody := `{
			"email": "user@example.com",
			"code": "123456",
			"new_password": "newpassword123"
		}`
		req, _ := http.NewRequest("POST", "/users/reset/confirm", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusTooManyRequests, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "rate_limited", response.Error)

		// Check rate limit headers
		assert.Equal(t, "10", w.Header().Get("X-RateLimit-Limit"))
		assert.Equal(t, "0", w.Header().Get("X-RateLimit-Remaining"))
		assert.NotEmpty(t, w.Header().Get("Retry-After"))
	})

	t.Run("rate limiter error", func(t *testing.T) {
		svc := &mockPasswordResetService{}
		limiter := &mockRateLimiter{
			allowErr: errors.New("redis error"),
		}
		r := setupPasswordResetRouter(svc, limiter)

		reqBody := `{
			"email": "user@example.com",
			"code": "123456",
			"new_password": "newpassword123"
		}`
		req, _ := http.NewRequest("POST", "/users/reset/confirm", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "internal_error", response.Error)
	})

	t.Run("service error - generic", func(t *testing.T) {
		svc := &mockPasswordResetService{
			confirmResetErr: errors.New("invalid code"),
		}
		limiter := &mockRateLimiter{
			allowResult:    true,
			allowRemaining: 9,
			allowResetAt:   time.Now().Add(time.Hour),
		}
		r := setupPasswordResetRouter(svc, limiter)

		reqBody := `{
			"email": "user@example.com",
			"code": "123456",
			"new_password": "newpassword123"
		}`
		req, _ := http.NewRequest("POST", "/users/reset/confirm", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid_request", response.Error)
	})

	t.Run("service error - password validation", func(t *testing.T) {
		svc := &mockPasswordResetService{
			confirmResetErr: errors.New("password must be at least 8 characters"),
		}
		limiter := &mockRateLimiter{
			allowResult:    true,
			allowRemaining: 9,
			allowResetAt:   time.Now().Add(time.Hour),
		}
		r := setupPasswordResetRouter(svc, limiter)

		reqBody := `{
			"email": "user@example.com",
			"code": "123456",
			"new_password": "short"
		}`
		req, _ := http.NewRequest("POST", "/users/reset/confirm", strings.NewReader(reqBody))
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
