package shipping_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/apperror"
	"github.com/leoferamos/aroma-sense/internal/dto"
	"github.com/leoferamos/aroma-sense/internal/handler/shipping"
)

type mockShippingService struct {
	options []dto.ShippingOption
	err     error
}

func (m mockShippingService) CalculateOptions(ctx context.Context, userID string, postalCode string) ([]dto.ShippingOption, error) {
	return m.options, m.err
}

func setupRouterWithUser(h *shipping.ShippingHandler) *gin.Engine {
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

func TestGetShippingOptions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		postalCode string
		svc        mockShippingService
		expectCode int
		assertion  func(t *testing.T, body []byte)
	}{
		{
			name:       "success",
			postalCode: "01234-567",
			svc:        mockShippingService{options: []dto.ShippingOption{{Carrier: "Correios", ServiceCode: "SEDEX", Price: 24.9, EstimatedDays: 2}}},
			expectCode: http.StatusOK,
			assertion: func(t *testing.T, body []byte) {
				var got []dto.ShippingOption
				if err := json.Unmarshal(body, &got); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				if len(got) != 1 || got[0].ServiceCode != "SEDEX" {
					t.Fatalf("unexpected response: %+v", got)
				}
			},
		},
		{
			name:       "invalid postal code",
			postalCode: "abc",
			svc:        mockShippingService{err: apperror.NewCodeMessage("invalid_postal_code", "invalid destination postal code")},
			expectCode: http.StatusBadRequest,
			assertion: func(t *testing.T, body []byte) {
				var errResp dto.ErrorResponse
				if err := json.Unmarshal(body, &errResp); err != nil {
					t.Fatalf("failed to unmarshal error response: %v", err)
				}
				if errResp.Error != "invalid_postal_code" {
					t.Fatalf("unexpected error message: %s", errResp.Error)
				}
			},
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			h := shipping.NewShippingHandler(tc.svc)
			r := setupRouterWithUser(h)

			req := httptest.NewRequest(http.MethodGet, "/shipping/options?postal_code="+tc.postalCode, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tc.expectCode {
				t.Fatalf("expected status %d, got %d", tc.expectCode, w.Code)
			}
			if tc.assertion != nil {
				tc.assertion(t, w.Body.Bytes())
			}
		})
	}
}
