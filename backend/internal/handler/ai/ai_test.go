package ai_test

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
	"github.com/leoferamos/aroma-sense/internal/handler/ai"
	"github.com/leoferamos/aroma-sense/internal/rate"
	"github.com/stretchr/testify/assert"
)

type AIServiceInterface interface {
	Recommend(ctx context.Context, message string, limit int) ([]dto.RecommendSuggestion, string, error)
}

type mockAIService struct {
	recommendResult []dto.RecommendSuggestion
	recommendReason string
	recommendErr    error
}

func (m *mockAIService) Recommend(ctx context.Context, message string, limit int) ([]dto.RecommendSuggestion, string, error) {
	return m.recommendResult, m.recommendReason, m.recommendErr
}

type mockRateLimiter struct {
	allowResult bool
	allowErr    error
	resetAt     time.Time
}

func (m *mockRateLimiter) Allow(ctx context.Context, bucket string, limit int, window time.Duration) (bool, int, time.Time, error) {
	return m.allowResult, 10, m.resetAt, m.allowErr
}

func setupAIHandler(svc AIServiceInterface, limiter rate.RateLimiter) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	handler := ai.NewAIHandler(svc, limiter)
	r.POST("/ai/recommend", handler.Recommend)
	return r
}

func TestAIHandler_Recommend(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockAIService{
			recommendResult: []dto.RecommendSuggestion{
				{ID: 1, Name: "Perfume 1", Slug: "perfume-1", Price: 29.99, Reason: "floral scent"},
			},
			recommendReason: "based on preferences",
		}
		limiter := &mockRateLimiter{allowResult: true}
		r := setupAIHandler(svc, limiter)

		reqBody := `{"message": "I like floral perfumes", "limit": 5}`
		req, _ := http.NewRequest("POST", "/ai/recommend", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.RecommendResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Suggestions, 1)
		assert.Equal(t, "based on preferences", response.Reasoning)
	})

	t.Run("invalid json", func(t *testing.T) {
		svc := &mockAIService{}
		limiter := &mockRateLimiter{allowResult: true}
		r := setupAIHandler(svc, limiter)

		reqBody := `{"message": "", "limit": 5}`
		req, _ := http.NewRequest("POST", "/ai/recommend", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid_request", response.Error)
	})

	t.Run("empty message", func(t *testing.T) {
		svc := &mockAIService{}
		limiter := &mockRateLimiter{allowResult: true}
		r := setupAIHandler(svc, limiter)

		reqBody := `{"message": "   ", "limit": 5}`
		req, _ := http.NewRequest("POST", "/ai/recommend", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid_request", response.Error)
	})

	t.Run("topic restricted", func(t *testing.T) {
		svc := &mockAIService{}
		limiter := &mockRateLimiter{allowResult: true}
		r := setupAIHandler(svc, limiter)

		reqBody := `{"message": "I like pizza", "limit": 5}`
		req, _ := http.NewRequest("POST", "/ai/recommend", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "topic_restricted", response.Error)
	})

	t.Run("rate limited", func(t *testing.T) {
		svc := &mockAIService{}
		limiter := &mockRateLimiter{allowResult: false, resetAt: time.Now().Add(5 * time.Second)}
		r := setupAIHandler(svc, limiter)

		reqBody := `{"message": "I like floral perfumes", "limit": 5}`
		req, _ := http.NewRequest("POST", "/ai/recommend", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusTooManyRequests, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "rate_limited", response.Error)
	})

	t.Run("service error", func(t *testing.T) {
		svc := &mockAIService{
			recommendErr: errors.New("database error"),
		}
		limiter := &mockRateLimiter{allowResult: true}
		r := setupAIHandler(svc, limiter)

		reqBody := `{"message": "I like floral perfumes", "limit": 5}`
		req, _ := http.NewRequest("POST", "/ai/recommend", strings.NewReader(reqBody))
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
