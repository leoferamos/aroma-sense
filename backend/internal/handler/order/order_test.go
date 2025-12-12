package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/dto"
	orderservice "github.com/leoferamos/aroma-sense/internal/service/order"
	"github.com/stretchr/testify/assert"
)

type mockOrderService struct {
	createOrderFromCartResult *dto.OrderResponse
	createOrderFromCartErr    error
	adminListOrdersResult     *dto.AdminOrdersResponse
	adminListOrdersErr        error
	getOrdersByUserResult     []dto.OrderResponse
	getOrdersByUserErr        error
}

func (m *mockOrderService) CreateOrderFromCart(userID string, req *dto.CreateOrderFromCartRequest) (*dto.OrderResponse, error) {
	return m.createOrderFromCartResult, m.createOrderFromCartErr
}

func (m *mockOrderService) AdminListOrders(status *string, startDate *time.Time, endDate *time.Time, page int, perPage int) (*dto.AdminOrdersResponse, error) {
	return m.adminListOrdersResult, m.adminListOrdersErr
}

func (m *mockOrderService) GetOrdersByUser(userID string) ([]dto.OrderResponse, error) {
	return m.getOrdersByUserResult, m.getOrdersByUserErr
}

func setupOrderRouter(svc orderservice.OrderService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Add middleware to simulate authentication
	r.Use(func(c *gin.Context) {
		c.Set("userID", "user-123")
		c.Next()
	})

	handler := NewOrderHandler(svc)
	r.POST("/orders", handler.CreateOrderFromCart)
	r.GET("/orders", handler.ListUserOrders)
	r.GET("/admin/orders", handler.ListOrders)
	return r
}

func setupOrderRouterUnauthenticated(svc orderservice.OrderService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	handler := NewOrderHandler(svc)
	r.POST("/orders", handler.CreateOrderFromCart)
	r.GET("/orders", handler.ListUserOrders)
	r.GET("/admin/orders", handler.ListOrders)
	return r
}

func TestOrderHandler_CreateOrderFromCart(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		orderResp := &dto.OrderResponse{PublicID: "order-123", Status: "pending", TotalAmount: 100.0}
		svc := &mockOrderService{
			createOrderFromCartResult: orderResp,
		}
		r := setupOrderRouter(svc)

		reqBody := `{
			"shipping_address": "Rua Teste, 123, São Paulo - SP, 01234-567",
			"payment_method": "credit_card"
		}`
		req, _ := http.NewRequest("POST", "/orders", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response dto.OrderResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "order-123", response.PublicID)
	})

	t.Run("invalid json", func(t *testing.T) {
		svc := &mockOrderService{}
		r := setupOrderRouter(svc)

		reqBody := `{"invalid": json}`
		req, _ := http.NewRequest("POST", "/orders", strings.NewReader(reqBody))
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
		svc := &mockOrderService{}
		gin.SetMode(gin.TestMode)
		r := gin.New()

		handler := NewOrderHandler(svc)
		r.POST("/orders", handler.CreateOrderFromCart)

		reqBody := `{
			"shipping_address": "Rua Teste, 123, São Paulo - SP, 01234-567",
			"payment_method": "credit_card"
		}`
		req, _ := http.NewRequest("POST", "/orders", strings.NewReader(reqBody))
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
		svc := &mockOrderService{
			createOrderFromCartErr: errors.New("database error"),
		}
		r := setupOrderRouter(svc)

		reqBody := `{
			"shipping_address": "Rua Teste, 123, São Paulo - SP, 01234-567",
			"payment_method": "credit_card"
		}`
		req, _ := http.NewRequest("POST", "/orders", strings.NewReader(reqBody))
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

func TestOrderHandler_ListUserOrders(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		orders := []dto.OrderResponse{
			{PublicID: "order-1", Status: "delivered", TotalAmount: 50.0},
			{PublicID: "order-2", Status: "pending", TotalAmount: 75.0},
		}
		svc := &mockOrderService{
			getOrdersByUserResult: orders,
		}
		r := setupOrderRouter(svc)

		req, _ := http.NewRequest("GET", "/orders", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []dto.OrderResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 2)
	})

	t.Run("unauthenticated", func(t *testing.T) {
		svc := &mockOrderService{}
		gin.SetMode(gin.TestMode)
		r := gin.New()

		handler := NewOrderHandler(svc)
		r.GET("/orders", handler.ListUserOrders)

		req, _ := http.NewRequest("GET", "/orders", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "unauthenticated", response.Error)
	})

	t.Run("service error", func(t *testing.T) {
		svc := &mockOrderService{
			getOrdersByUserErr: errors.New("database error"),
		}
		r := setupOrderRouter(svc)

		req, _ := http.NewRequest("GET", "/orders", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "internal_error", response.Error)
	})
}

func TestOrderHandler_ListOrders(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		adminResp := &dto.AdminOrdersResponse{
			Orders: []dto.AdminOrderItem{
				{ID: 1, PublicID: "order-123", UserID: "user-456", TotalAmount: 100.0, Status: "pending", CreatedAt: time.Now()},
			},
			Meta: struct {
				Pagination dto.PaginationMeta `json:"pagination"`
				Stats      dto.StatsMeta      `json:"stats"`
			}{
				Pagination: dto.PaginationMeta{Page: 1, PerPage: 25, TotalPages: 1, TotalCount: 1},
				Stats:      dto.StatsMeta{TotalRevenue: 100.0, AverageOrderValue: 100.0},
			},
		}
		svc := &mockOrderService{
			adminListOrdersResult: adminResp,
		}
		r := setupOrderRouter(svc)

		req, _ := http.NewRequest("GET", "/admin/orders", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.AdminOrdersResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Orders, 1)
		assert.Equal(t, 1, response.Meta.Pagination.TotalCount)
	})

	t.Run("invalid status", func(t *testing.T) {
		svc := &mockOrderService{}
		r := setupOrderRouter(svc)

		req, _ := http.NewRequest("GET", "/admin/orders?status=invalid", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid_request", response.Error)
	})

	t.Run("invalid date", func(t *testing.T) {
		svc := &mockOrderService{}
		r := setupOrderRouter(svc)

		req, _ := http.NewRequest("GET", "/admin/orders?start_date=invalid", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid_request", response.Error)
	})

	t.Run("service error", func(t *testing.T) {
		svc := &mockOrderService{
			adminListOrdersErr: errors.New("database error"),
		}
		r := setupOrderRouter(svc)

		req, _ := http.NewRequest("GET", "/admin/orders", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "internal_error", response.Error)
	})
}
