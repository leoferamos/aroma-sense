package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"context"

	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/dto"
	"github.com/leoferamos/aroma-sense/internal/service"
)

type mockShippingService struct {
	options []dto.ShippingOption
	err     error
}

func (m mockShippingService) CalculateOptions(ctx context.Context, userID string, postalCode string) ([]dto.ShippingOption, error) {
	return m.options, m.err
}

func setupRouterWithUser(h *ShippingHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	// Inject userID into context as if authenticated
	r.Use(func(c *gin.Context) {
		c.Set("userID", "user-123")
		c.Next()
	})
	r.GET("/shipping/options", h.GetShippingOptions)
	return r
}

func TestGetShippingOptions_Success(t *testing.T) {
	mockSvc := mockShippingService{
		options: []dto.ShippingOption{{Carrier: "Correios", ServiceCode: "SEDEX", Price: 24.9, EstimatedDays: 2}},
		err:     nil,
	}
	h := NewShippingHandler(mockSvc)
	r := setupRouterWithUser(h)

	req := httptest.NewRequest(http.MethodGet, "/shipping/options?postal_code=01234-567", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	var got []dto.ShippingOption
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if len(got) != 1 || got[0].ServiceCode != "SEDEX" {
		t.Fatalf("unexpected response: %+v", got)
	}
}

func TestGetShippingOptions_ErrorMapping(t *testing.T) {
	mockSvc := mockShippingService{
		options: nil,
		err:     service.ErrInvalidPostalCode,
	}
	h := NewShippingHandler(mockSvc)
	r := setupRouterWithUser(h)

	req := httptest.NewRequest(http.MethodGet, "/shipping/options?postal_code=abc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
	var errResp dto.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("failed to unmarshal error response: %v", err)
	}
	if errResp.Error != "invalid postal_code" {
		t.Fatalf("unexpected error message: %s", errResp.Error)
	}
}
