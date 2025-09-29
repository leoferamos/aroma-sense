package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/dto"
	"github.com/leoferamos/aroma-sense/internal/service"
)

type CartHandler struct {
	cartService service.CartService
}

func NewCartHandler(cartService service.CartService) *CartHandler {
	return &CartHandler{cartService: cartService}
}

// GetCart retrieves the current user's cart
//
// @Summary      Get current user's cart
// @Description  Retrieves the shopping cart for the authenticated user with items, quantities and totals
// @Tags         cart
// @Accept       json
// @Produce      json
// @Param        Authorization  header  string  true  "Bearer JWT token"
// @Success      200  {object}  dto.CartResponse  "User's cart with items and totals"
// @Failure      401  {object}  dto.ErrorResponse  "Unauthorized"
// @Failure      404  {object}  dto.ErrorResponse  "Cart not found"
// @Failure      500  {object}  dto.ErrorResponse  "Internal server error"
// @Router       /cart [get]
// @Security     BearerAuth
func (h *CartHandler) GetCart(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "user not authenticated"})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "invalid user ID"})
		return
	}

	// Get user's cart
	cartResponse, err := h.cartService.GetCartResponse(userIDStr)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "cart not found"})
		return
	}

	c.JSON(http.StatusOK, cartResponse)
}

// AddItem adds an item to the user's cart
//
// @Summary      Add item to cart
// @Description  Adds a product to the user's shopping cart. If item already exists, increases quantity.
// @Tags         cart
// @Accept       json
// @Produce      json
// @Param        Authorization  header  string                true   "Bearer JWT token"
// @Param        request        body    dto.AddToCartRequest  true   "Product ID and quantity to add"
// @Success      200  {object}  dto.CartResponse    "Updated cart with new item"
// @Failure      400  {object}  dto.ErrorResponse   "Invalid request body"
// @Failure      401  {object}  dto.ErrorResponse   "Unauthorized"
// @Failure      500  {object}  dto.ErrorResponse   "Internal server error"
// @Router       /cart [post]
// @Security     BearerAuth
func (h *CartHandler) AddItem(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "user not authenticated"})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "invalid user ID"})
		return
	}

	var req dto.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body"})
		return
	}

	// Add item to cart
	cartResponse, err := h.cartService.AddItemToCart(userIDStr, req.ProductID, req.Quantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, cartResponse)
}
