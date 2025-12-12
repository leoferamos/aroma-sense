package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/dto"
	paymentservice "github.com/leoferamos/aroma-sense/internal/service/payment"
	"github.com/stretchr/testify/assert"
)

type mockPaymentService struct {
	createIntentResult  *paymentservice.PaymentIntentResult
	createIntentErr     error
	handleWebhookResult *paymentservice.PaymentWebhookPayload
	handleWebhookErr    error
}

func (m *mockPaymentService) CreateIntent(ctx context.Context, userID string, req *dto.CreatePaymentIntentRequest) (*paymentservice.PaymentIntentResult, error) {
	return m.createIntentResult, m.createIntentErr
}

func (m *mockPaymentService) HandleWebhook(ctx context.Context, payload []byte, signature string) (*paymentservice.PaymentWebhookPayload, error) {
	return m.handleWebhookResult, m.handleWebhookErr
}

func setupPaymentRouter(svc paymentservice.PaymentService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Add middleware to simulate authentication
	r.Use(func(c *gin.Context) {
		c.Set("userID", "user-123")
		c.Next()
	})

	handler := NewPaymentHandler(svc)
	r.POST("/payments/intent", handler.CreateIntent)
	r.POST("/payments/webhook", handler.HandleWebhook)
	return r
}

func TestPaymentHandler_CreateIntent(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		intentResult := &paymentservice.PaymentIntentResult{
			ID:           "pi_123",
			ClientSecret: "secret_123",
		}
		svc := &mockPaymentService{
			createIntentResult: intentResult,
		}
		r := setupPaymentRouter(svc)

		reqBody := `{
			"shipping_address": "Rua Teste, 123, São Paulo - SP, 01234-567"
		}`
		req, _ := http.NewRequest("POST", "/payments/intent", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "pi_123", response["payment_intent_id"])
		assert.Equal(t, "secret_123", response["client_secret"])
	})

	t.Run("invalid json", func(t *testing.T) {
		svc := &mockPaymentService{}
		r := setupPaymentRouter(svc)

		reqBody := `{"invalid": json}`
		req, _ := http.NewRequest("POST", "/payments/intent", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid_request", response.Error)
	})

	t.Run("unauthenticated", func(t *testing.T) {
		svc := &mockPaymentService{}
		gin.SetMode(gin.TestMode)
		r := gin.New()

		handler := NewPaymentHandler(svc)
		r.POST("/payments/intent", handler.CreateIntent)

		reqBody := `{
			"shipping_address": "Rua Teste, 123, São Paulo - SP, 01234-567"
		}`
		req, _ := http.NewRequest("POST", "/payments/intent", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "unauthenticated", response.Error)
	})

	t.Run("service error", func(t *testing.T) {
		svc := &mockPaymentService{
			createIntentErr: errors.New("payment error"),
		}
		r := setupPaymentRouter(svc)

		reqBody := `{
			"shipping_address": "Rua Teste, 123, São Paulo - SP, 01234-567"
		}`
		req, _ := http.NewRequest("POST", "/payments/intent", strings.NewReader(reqBody))
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

func TestPaymentHandler_HandleWebhook(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		webhookPayload := &paymentservice.PaymentWebhookPayload{
			IntentID:      "pi_123",
			Status:        "succeeded",
			Amount:        1000,
			Currency:      "usd",
			CustomerEmail: "test@example.com",
			Metadata:      map[string]string{"order_id": "order-123"},
		}
		svc := &mockPaymentService{
			handleWebhookResult: webhookPayload,
		}
		r := setupPaymentRouter(svc)

		reqBody := `{"type": "payment_intent.succeeded"}`
		req, _ := http.NewRequest("POST", "/payments/webhook", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Stripe-Signature", "signature123")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("missing signature", func(t *testing.T) {
		svc := &mockPaymentService{}
		r := setupPaymentRouter(svc)

		reqBody := `{"type": "payment_intent.succeeded"}`
		req, _ := http.NewRequest("POST", "/payments/webhook", bytes.NewReader([]byte(reqBody)))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "missing_signature", response.Error)
	})

	t.Run("service error", func(t *testing.T) {
		svc := &mockPaymentService{
			handleWebhookErr: errors.New("webhook error"),
		}
		r := setupPaymentRouter(svc)

		reqBody := `{"type": "payment_intent.succeeded"}`
		req, _ := http.NewRequest("POST", "/payments/webhook", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Stripe-Signature", "signature123")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "internal_error", response.Error)
	})
}
