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

type ChatHandler struct {
	chat  *service.ChatService
	limit rate.RateLimiter
}

func NewChatHandler(chat *service.ChatService, limiter rate.RateLimiter) *ChatHandler {
	return &ChatHandler{chat: chat, limit: limiter}
}

// Chat handles a conversational LLM turn.
func (h *ChatHandler) Chat(c *gin.Context) {
	var req dto.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Message) == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "mensagem inválida"})
		return
	}
	ip := clientIP(c)
	if allowed, _, reset, err := h.limit.Allow(c.Request.Context(), "chat:"+ip, 20, time.Minute); err != nil || !allowed {
		_ = reset
		c.JSON(http.StatusTooManyRequests, dto.ErrorResponse{Error: "limite excedido, aguarde alguns segundos"})
		return
	}
	resp, err := h.chat.Chat(c.Request.Context(), req.SessionID, req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "falha ao gerar resposta"})
		return
	}
	c.JSON(http.StatusOK, resp)
}
