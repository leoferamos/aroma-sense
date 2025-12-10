package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/dto"
	handlererrors "github.com/leoferamos/aroma-sense/internal/handler/errors"
	"github.com/leoferamos/aroma-sense/internal/rate"
	chatservice "github.com/leoferamos/aroma-sense/internal/service/chat"
)

type AIHandler struct {
	svc   *chatservice.AIService
	limit rate.RateLimiter
}

func NewAIHandler(svc *chatservice.AIService, limiter rate.RateLimiter) *AIHandler {
	return &AIHandler{svc: svc, limit: limiter}
}

// Recommend returns AI-powered product suggestions based on user preferences.
//
// @Summary      Get AI product recommendations
// @Description  Returns personalized product recommendations based on user's fragrance preferences and message
// @Tags         ai
// @Accept       json
// @Produce      json
// @Param        request  body  dto.RecommendRequest  true  "Recommendation request with user message and limit"
// @Success      200  {object}  dto.RecommendResponse  "Product recommendations with reasoning"
// @Failure      400  {object}  dto.ErrorResponse     "Error code: invalid_request or topic_restricted"
// @Failure      429  {object}  dto.ErrorResponse     "Error code: rate_limited"
// @Failure      500  {object}  dto.ErrorResponse     "Error code: internal_error"
// @Router       /ai/recommend [post]
func (h *AIHandler) Recommend(c *gin.Context) {
	var req dto.RecommendRequest
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Message) == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request"})
		return
	}

	// Basic topic guard: keep conversation perfume-related.
	lower := strings.ToLower(req.Message)
	if !containsAny(lower, []string{"perfume", "fragr", "cheiro", "aroma", "odor", "eau"}) {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "topic_restricted"})
		return
	}

	// Rate limit per IP for safety (10 req/min)
	ip := clientIP(c)
	if allowed, _, reset, err := h.limit.Allow(c.Request.Context(), "ai:"+ip, 10, time.Minute); err != nil || !allowed {
		retry := int(time.Until(reset).Seconds())
		if retry < 1 {
			retry = 5
		}
		c.Header("Retry-After", "5")
		c.JSON(http.StatusTooManyRequests, dto.ErrorResponse{Error: "rate_limited"})
		return
	}

	suggestions, reason, err := h.svc.Recommend(c.Request.Context(), req.Message, req.Limit)
	if err != nil {
		if status, code, ok := handlererrors.MapServiceError(err); ok {
			c.JSON(status, dto.ErrorResponse{Error: code})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error"})
		return
	}

	c.JSON(http.StatusOK, dto.RecommendResponse{Suggestions: suggestions, Reasoning: reason})
}

func containsAny(s string, subs []string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

func clientIP(c *gin.Context) string {
	// honor X-Forwarded-For if present
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	return c.ClientIP()
}
