package chat

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

type ChatHandler struct {
	chat  *chatservice.ChatService
	limit rate.RateLimiter
}

func NewChatHandler(chat *chatservice.ChatService, limiter rate.RateLimiter) *ChatHandler {
	return &ChatHandler{chat: chat, limit: limiter}
}

// Chat handles a conversational LLM turn.
//
// @Summary      Chat with AI assistant
// @Description  Send a message to the AI assistant and receive a conversational response with product suggestions
// @Tags         ai
// @Accept       json
// @Produce      json
// @Param        request  body  dto.ChatRequest  true  "Chat message with optional session ID"
// @Success      200  {object}  dto.ChatResponse  "AI response with reply and product suggestions"
// @Failure      400  {object}  dto.ErrorResponse "Error code: invalid_request"
// @Failure      429  {object}  dto.ErrorResponse "Error code: rate_limited"
// @Failure      500  {object}  dto.ErrorResponse "Error code: internal_error"
// @Router       /ai/chat [post]
func (h *ChatHandler) Chat(c *gin.Context) {
	var req dto.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Message) == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request"})
		return
	}
	ip := clientIP(c)
	if allowed, _, reset, err := h.limit.Allow(c.Request.Context(), "chat:"+ip, 20, time.Minute); err != nil || !allowed {
		_ = reset
		c.JSON(http.StatusTooManyRequests, dto.ErrorResponse{Error: "rate_limited"})
		return
	}
	resp, err := h.chat.Chat(c.Request.Context(), req.SessionID, req.Message)
	if err != nil {
		if status, code, ok := handlererrors.MapServiceError(err); ok {
			c.JSON(status, dto.ErrorResponse{Error: code})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error"})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func clientIP(c *gin.Context) string {
	// honor X-Forwarded-For if present
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	return c.ClientIP()
}
