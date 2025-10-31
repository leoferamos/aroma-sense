package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/dto"
	"github.com/leoferamos/aroma-sense/internal/model"
	"github.com/leoferamos/aroma-sense/internal/service"
)

type OrderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{orderService}
}

const maxPerPage = 100

// CreateOrderFromCart handles order creation from the user's cart
//
// @Summary      Create a new order
// @Description  Creates a new order from the user's cart, validates stock, deducts products, and clears the cart
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        order  body  dto.CreateOrderFromCartRequest  true  "Order data (shipping address, payment method)"
// @Success      201  {object}  dto.OrderResponse      "Order created successfully"
// @Failure      400  {object}  dto.ErrorResponse      "Invalid request data or empty cart"
// @Failure      401  {object}  dto.ErrorResponse      "Unauthorized"
// @Router       /orders [post]
// @Security     BearerAuth
func (h *OrderHandler) CreateOrderFromCart(c *gin.Context) {
	var req dto.CreateOrderFromCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid request"})
		return
	}

	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "Unauthorized"})
		return
	}

	orderResp, err := h.orderService.CreateOrderFromCart(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, orderResp)
}

// ListOrders allows admin to list orders with filters, pagination and stats
// @Summary      List all orders
// @Description  Returns paginated list of orders with optional filters (status, date range)
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        status     query    string  false  "Filter by status"
// @Param        start_date query    string  false  "Start date (YYYY-MM-DD)"
// @Param        end_date   query    string  false  "End date (YYYY-MM-DD)"
// @Param        page       query    int     false  "Page number (1-based)"
// @Param        per_page   query    int     false  "Items per page"
// @Success      200  {object}  dto.AdminOrdersResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Router       /admin/orders [get]
// @Security     BearerAuth
func (h *OrderHandler) ListOrders(c *gin.Context) {

	// Parse query params
	statusParam := c.Query("status")
	var status *string
	if statusParam != "" {
		// Validate status against known OrderStatus values
		// If invalid, return 400 to the client
		switch model.OrderStatus(statusParam) {
		case model.OrderStatusPending, model.OrderStatusProcessing, model.OrderStatusShipped, model.OrderStatusDelivered, model.OrderStatusCancelled:
			// valid
		default:
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid status"})
			return
		}
		status = &statusParam
	}

	layout := "2006-01-02"
	var startDatePtr *time.Time
	var endDatePtr *time.Time
	if s := c.Query("start_date"); s != "" {
		if t, err := time.Parse(layout, s); err == nil {
			startDatePtr = &t
		} else {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid start_date format, use YYYY-MM-DD"})
			return
		}
	}
	if s := c.Query("end_date"); s != "" {
		if t, err := time.Parse(layout, s); err == nil {
			end := t.Add(24*time.Hour - time.Nanosecond)
			endDatePtr = &end
		} else {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid end_date format, use YYYY-MM-DD"})
			return
		}
	}

	page := 1
	perPage := 25
	if p := c.Query("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		} else {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid page"})
			return
		}
	}
	if pp := c.Query("per_page"); pp != "" {
		if v, err := strconv.Atoi(pp); err == nil && v > 0 {
			perPage = v
		} else {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid per_page"})
			return
		}
	}

	// Enforce maxPerPage to avoid heavy responses
	if perPage > maxPerPage {
		perPage = maxPerPage
	}

	resp, err := h.orderService.AdminListOrders(status, startDatePtr, endDatePtr, page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to list orders"})
		return
	}
	c.JSON(http.StatusOK, resp)
}
