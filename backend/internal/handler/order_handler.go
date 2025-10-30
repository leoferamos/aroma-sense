package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/dto"
	"github.com/leoferamos/aroma-sense/internal/service"
)

type OrderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{orderService}
}

// CreateOrderFromCart handles order creation from the user's cart
//
// @Summary      Create a new order
// @Description  Creates a new order from the user's cart, validates stock, deducts products, and clears the cart
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        Authorization  header  string  true  "Bearer JWT token"
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

	userID := c.GetString("user_id")
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
