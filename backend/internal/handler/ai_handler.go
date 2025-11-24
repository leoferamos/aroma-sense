package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/dto"
	"github.com/leoferamos/aroma-sense/internal/rate"
	"github.com/leoferamos/aroma-sense/internal/service"
)

type AIHandler struct {
	svc   *service.AIService
	limit rate.RateLimiter
}

func NewAIHandler(svc *service.AIService, limiter rate.RateLimiter) *AIHandler {
	return &AIHandler{svc: svc, limit: limiter}
}

// Recommend returns product suggestions.
func (h *AIHandler) Recommend(c *gin.Context) {
	var req dto.RecommendRequest
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Message) == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "mensagem inválida"})
		return
	}

	// Basic topic guard: we keep it perfume-related to avoid abuse.
	lower := strings.ToLower(req.Message)
	if !containsAny(lower, []string{"perfume", "fragr", "cheiro", "aroma", "odor", "eau"}) {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Vamos falar de perfumes e fragrâncias. Conte suas preferências :)"})
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
		c.JSON(http.StatusTooManyRequests, dto.ErrorResponse{Error: "muitas requisições, tente novamente em alguns segundos"})
		return
	}

	suggestions, reason, err := h.svc.Recommend(c.Request.Context(), req.Message, req.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "não foi possível gerar recomendações"})
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
