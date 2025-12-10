package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/apperror"
	"github.com/leoferamos/aroma-sense/internal/dto"
	handler "github.com/leoferamos/aroma-sense/internal/handler"
	paymentservice "github.com/leoferamos/aroma-sense/internal/service/payment"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockPaymentService struct{ mock.Mock }

func (m *mockPaymentService) CreateIntent(ctx context.Context, userID string, req *dto.CreatePaymentIntentRequest) (*paymentservice.PaymentIntentResult, error) {
	args := m.Called(ctx, userID, req)
	if res, ok := args.Get(0).(*paymentservice.PaymentIntentResult); ok {
		return res, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockPaymentService) HandleWebhook(ctx context.Context, payload []byte, signature string) (*paymentservice.PaymentWebhookPayload, error) {
	args := m.Called(ctx, payload, signature)
	if res, ok := args.Get(0).(*paymentservice.PaymentWebhookPayload); ok {
		return res, args.Error(1)
	}
	return nil, args.Error(1)
}

func setupPaymentRouter(svc *mockPaymentService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	h := handler.NewPaymentHandler(svc)
	r := gin.New()
	r.POST("/payments/intent", func(c *gin.Context) {
		c.Set("userID", "user-123")
		h.CreateIntent(c)
	})
	r.POST("/payments/webhook", h.HandleWebhook)
	return r
}

func TestPaymentHandler_CreateIntent(t *testing.T) {
	t.Parallel()

	r := gin.New()
	svc := new(mockPaymentService)
	h := handler.NewPaymentHandler(svc)
	r.POST("/payments/intent", func(c *gin.Context) {
		c.Set("userID", "user-123")
		h.CreateIntent(c)
	})

	payload := dto.CreatePaymentIntentRequest{ShippingAddress: "12345-678", OrderPublicID: "order-1"}
	body, err := json.Marshal(payload)
	require.NoError(t, err)

	svc.On("CreateIntent", mock.Anything, "user-123", mock.AnythingOfType("*dto.CreatePaymentIntentRequest")).
		Return(&paymentservice.PaymentIntentResult{ID: "pi_123", ClientSecret: "secret"}, nil).
		Once()

	req := httptest.NewRequest(http.MethodPost, "/payments/intent", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	require.Equal(t, http.StatusOK, res.Code)

	var respBody map[string]string
	require.NoError(t, json.Unmarshal(res.Body.Bytes(), &respBody))
	assert.Equal(t, "pi_123", respBody["payment_intent_id"])
	assert.Equal(t, "secret", respBody["client_secret"])

	svc.AssertExpectations(t)

	t.Run("invalid JSON returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/payments/intent", bytes.NewBufferString("{"))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	t.Run("unauthenticated returns 401", func(t *testing.T) {
		rNoAuth := gin.New()
		rNoAuth.POST("/payments/intent", h.CreateIntent)

		req := httptest.NewRequest(http.MethodPost, "/payments/intent", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		rNoAuth.ServeHTTP(res, req)
		assert.Equal(t, http.StatusUnauthorized, res.Code)
	})

	t.Run("mapped service error", func(t *testing.T) {
		svc.On("CreateIntent", mock.Anything, "user-123", &payload).
			Return(nil, apperror.NewCodeMessage("invalid_request", "")).Once()

		req := httptest.NewRequest(http.MethodPost, "/payments/intent", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Code)
		svc.AssertExpectations(t)
	})

	t.Run("generic service error", func(t *testing.T) {
		svc.On("CreateIntent", mock.Anything, "user-123", &payload).
			Return(nil, assert.AnError).Once()

		req := httptest.NewRequest(http.MethodPost, "/payments/intent", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		svc.AssertExpectations(t)
	})
}

func TestPaymentHandler_HandleWebhook(t *testing.T) {
	t.Parallel()

	svc := new(mockPaymentService)
	r := setupPaymentRouter(svc)

	goodBody := []byte(`{"id":"evt_1"}`)

	t.Run("missing signature returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/payments/webhook", bytes.NewReader(goodBody))
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	t.Run("mapped service error", func(t *testing.T) {
		svc.On("HandleWebhook", mock.Anything, goodBody, "sig").
			Return(nil, apperror.NewCodeMessage("invalid_webhook", "")).Once()

		req := httptest.NewRequest(http.MethodPost, "/payments/webhook", bytes.NewReader(goodBody))
		req.Header.Set("Stripe-Signature", "sig")
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Code)
		svc.AssertExpectations(t)
	})

	t.Run("generic service error", func(t *testing.T) {
		svc.On("HandleWebhook", mock.Anything, goodBody, "sig").
			Return(nil, assert.AnError).Once()

		req := httptest.NewRequest(http.MethodPost, "/payments/webhook", bytes.NewReader(goodBody))
		req.Header.Set("Stripe-Signature", "sig")
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		svc.AssertExpectations(t)
	})

	t.Run("success returns 200", func(t *testing.T) {
		svc.On("HandleWebhook", mock.Anything, goodBody, "sig").
			Return(&paymentservice.PaymentWebhookPayload{}, nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/payments/webhook", bytes.NewReader(goodBody))
		req.Header.Set("Stripe-Signature", "sig")
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Code)
		svc.AssertExpectations(t)
	})
}
