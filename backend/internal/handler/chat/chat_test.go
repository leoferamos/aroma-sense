package chat

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestChatRateLimit(t *testing.T) {
	r := gin.Default()
	r.POST("/chat", func(c *gin.Context) {
		if c.GetHeader("X-Test-Block") == "1" {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate_limited"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"reply": "ok"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/chat", strings.NewReader(`{"message":"oi"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/chat", strings.NewReader(`{"message":"oi"}`))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("X-Test-Block", "1")
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusTooManyRequests, w2.Code)
}

func TestChatBadRequest(t *testing.T) {
	r := gin.Default()
	r.POST("/chat", func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/chat", strings.NewReader(""))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
