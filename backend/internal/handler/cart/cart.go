package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/dto"
	handlererrors "github.com/leoferamos/aroma-sense/internal/handler/errors"
	cartservice "github.com/leoferamos/aroma-sense/internal/service/cart"
)

type CartHandler struct {
	cartService cartservice.CartService
}

func NewCartHandler(cartService cartservice.CartService) *CartHandler {
	return &CartHandler{cartService: cartService}
}

// GetCart retrieves the current user's cart
//
// @Summary      Get current user's cart
// @Description  Retrieves the shopping cart for the authenticated user with items, quantities and totals
// @Tags         cart
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.CartResponse  "User's cart with items and totals"
// @Failure      401  {object}  dto.ErrorResponse  "Error code: unauthenticated"
// @Failure      404  {object}  dto.ErrorResponse  "Error code: cart_not_found"
// @Failure      500  {object}  dto.ErrorResponse  "Error code: internal_error"
// @Router       /cart [get]
// @Security     BearerAuth
func (h *CartHandler) GetCart(c *gin.Context) {

	userIDStr := c.GetString("userID")
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthenticated"})
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
// @Param        request        body    dto.AddToCartRequest  true   "Product slug and quantity to add"
// @Success      200  {object}  dto.CartResponse    "Updated cart with new item"
// @Failure      400  {object}  dto.ErrorResponse   "Error code: invalid_request"
// @Failure      401  {object}  dto.ErrorResponse   "Error code: unauthenticated"
// @Failure      404  {object}  dto.ErrorResponse   "Error code: product_not_found"
// @Failure      409  {object}  dto.ErrorResponse   "Error code: insufficient_stock"
// @Failure      500  {object}  dto.ErrorResponse   "Error code: internal_error"
// @Router       /cart [post]
// @Security     BearerAuth
func (h *CartHandler) AddItem(c *gin.Context) {

	userIDStr := c.GetString("userID")
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthenticated"})
		return
	}

	var req dto.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request"})
		return
	}

	// Add item to cart
	cartResponse, err := h.cartService.AddItemToCart(userIDStr, req.ProductSlug, req.Quantity)
	if err != nil {
		if status, code, ok := handlererrors.MapServiceError(err); ok {
			c.JSON(status, dto.ErrorResponse{Error: code})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error"})
		return
	}

	c.JSON(http.StatusOK, cartResponse)
}

// UpdateItemQuantity updates the quantity of a specific cart item
//
// @Summary      Update cart item quantity
// @Description  Updates the quantity of a specific item in the user's cart. If quantity is 0, removes the item.
// @Tags         cart
// @Accept       json
// @Produce      json
// @Param        productSlug    path    string                       true   "Product slug"
// @Param        request        body    dto.UpdateCartItemRequest    true   "New quantity (0 to remove item)"
// @Success      200  {object}  dto.CartResponse    "Updated cart"
// @Failure      400  {object}  dto.ErrorResponse   "Error code: invalid_request"
// @Failure      401  {object}  dto.ErrorResponse   "Error code: unauthenticated"
// @Failure      404  {object}  dto.ErrorResponse   "Error code: cart_item_not_found"
// @Failure      409  {object}  dto.ErrorResponse   "Error code: insufficient_stock"
// @Failure      500  {object}  dto.ErrorResponse   "Error code: internal_error"
// @Router       /cart/items/{productSlug} [patch]
// @Security     BearerAuth
func (h *CartHandler) UpdateItemQuantity(c *gin.Context) {

	userIDStr := c.GetString("userID")
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthenticated"})
		return
	}

	// Get product slug from URL parameter
	productSlug := c.Param("productSlug")
	if productSlug == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request"})
		return
	}

	var req dto.UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request"})
		return
	}

	// Update item quantity
	cartResponse, err := h.cartService.UpdateItemQuantityBySlug(userIDStr, productSlug, req.Quantity)
	if err != nil {
		if status, code, ok := handlererrors.MapServiceError(err); ok {
			c.JSON(status, dto.ErrorResponse{Error: code})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error"})
		return
	}

	c.JSON(http.StatusOK, cartResponse)
}

// RemoveItem removes a specific item from the user's cart
//
// @Summary      Remove item from cart
// @Description  Removes a specific item from the user's cart completely
// @Tags         cart
// @Accept       json
// @Produce      json
// @Param        productSlug    path    string true   "Product slug"
// @Success      200  {object}  dto.CartResponse    "Updated cart after item removal"
// @Failure      400  {object}  dto.ErrorResponse   "Error code: invalid_request"
// @Failure      401  {object}  dto.ErrorResponse   "Error code: unauthenticated"
// @Failure      404  {object}  dto.ErrorResponse   "Error code: cart_item_not_found"
// @Failure      500  {object}  dto.ErrorResponse   "Error code: internal_error"
// @Router       /cart/items/{productSlug} [delete]
// @Security     BearerAuth
func (h *CartHandler) RemoveItem(c *gin.Context) {

	userIDStr := c.GetString("userID")
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthenticated"})
		return
	}

	// Get product slug from URL parameter
	productSlug := c.Param("productSlug")
	if productSlug == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request"})
		return
	}

	// Remove item from cart
	cartResponse, err := h.cartService.RemoveItemBySlug(userIDStr, productSlug)
	if err != nil {
		if status, code, ok := handlererrors.MapServiceError(err); ok {
			c.JSON(status, dto.ErrorResponse{Error: code})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error"})
		return
	}

	c.JSON(http.StatusOK, cartResponse)
}

// ClearCart removes all items from the user's cart
//
// @Summary      Clear cart
// @Description  Removes all items from the user's cart, returning an empty cart
// @Tags         cart
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.CartResponse    "Empty cart after clearing all items"
// @Failure      401  {object}  dto.ErrorResponse   "Error code: unauthenticated"
// @Failure      404  {object}  dto.ErrorResponse   "Error code: cart_not_found"
// @Failure      500  {object}  dto.ErrorResponse   "Error code: internal_error"
// @Router       /cart [delete]
// @Security     BearerAuth
func (h *CartHandler) ClearCart(c *gin.Context) {

	userIDStr := c.GetString("userID")
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthenticated"})
		return
	}

	// Clear all items from cart
	cartResponse, err := h.cartService.ClearCart(userIDStr)
	if err != nil {
		if status, code, ok := handlererrors.MapServiceError(err); ok {
			c.JSON(status, dto.ErrorResponse{Error: code})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error"})
		return
	}

	c.JSON(http.StatusOK, cartResponse)
}
