package handler

import (
	"net/http"
	"strconv"

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
// @Success      200  {object}  dto.CartResponse  "User's cart with items and totals"
// @Failure      401  {object}  dto.ErrorResponse  "Unauthorized"
// @Failure      404  {object}  dto.ErrorResponse  "Cart not found"
// @Failure      500  {object}  dto.ErrorResponse  "Internal server error"
// @Router       /cart [get]
// @Security     BearerAuth
func (h *CartHandler) GetCart(c *gin.Context) {

	userIDStr := c.GetString("userID")
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "user not authenticated"})
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
// @Failure      400  {object}  dto.ErrorResponse   "Invalid request body or insufficient stock"
// @Failure      401  {object}  dto.ErrorResponse   "Unauthorized"
// @Failure      404  {object}  dto.ErrorResponse   "Product not found"
// @Failure      409  {object}  dto.ErrorResponse   "Product out of stock"
// @Failure      500  {object}  dto.ErrorResponse   "Internal server error"
// @Router       /cart [post]
// @Security     BearerAuth
func (h *CartHandler) AddItem(c *gin.Context) {

	userIDStr := c.GetString("userID")
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "user not authenticated"})
		return
	}

	var req dto.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body"})
		return
	}

	// Add item to cart
	cartResponse, err := h.cartService.AddItemToCart(userIDStr, req.ProductSlug, req.Quantity)
	if err != nil {
		// Handle different types of errors with appropriate HTTP status codes
		if err.Error() == "product not found" {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
			return
		}
		if err.Error() == "product out of stock" {
			c.JSON(http.StatusConflict, dto.ErrorResponse{Error: err.Error()})
			return
		}
		// For insufficient stock or other validation errors
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
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
// @Param        itemId         path    int                          true   "Cart item ID"
// @Param        request        body    dto.UpdateCartItemRequest    true   "New quantity (0 to remove item)"
// @Success      200  {object}  dto.CartResponse    "Updated cart"
// @Failure      400  {object}  dto.ErrorResponse   "Invalid request body, item ID, or insufficient stock"
// @Failure      401  {object}  dto.ErrorResponse   "Unauthorized"
// @Failure      404  {object}  dto.ErrorResponse   "Cart item or product not found"
// @Failure      409  {object}  dto.ErrorResponse   "Product out of stock"
// @Failure      500  {object}  dto.ErrorResponse   "Internal server error"
// @Router       /cart/items/{itemId} [patch]
// @Security     BearerAuth
func (h *CartHandler) UpdateItemQuantity(c *gin.Context) {

	userIDStr := c.GetString("userID")
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "user not authenticated"})
		return
	}

	// Get item ID from URL parameter
	itemIDParam := c.Param("itemId")
	var itemID uint
	if id, err := strconv.ParseUint(itemIDParam, 10, 32); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid item ID"})
		return
	} else {
		itemID = uint(id)
	}

	var req dto.UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body"})
		return
	}

	// Update item quantity
	cartResponse, err := h.cartService.UpdateItemQuantity(userIDStr, itemID, req.Quantity)
	if err != nil {
		if err.Error() == "cart item not found" || err.Error() == "cart item not found in user's cart" {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
			return
		}
		if err.Error() == "product not found" {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
			return
		}
		if err.Error() == "product out of stock" {
			c.JSON(http.StatusConflict, dto.ErrorResponse{Error: err.Error()})
			return
		}
		// For insufficient stock or other validation errors
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
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
// @Param        itemId         path    int     true   "Cart item ID"
// @Success      200  {object}  dto.CartResponse    "Updated cart after item removal"
// @Failure      400  {object}  dto.ErrorResponse   "Invalid item ID"
// @Failure      401  {object}  dto.ErrorResponse   "Unauthorized"
// @Failure      404  {object}  dto.ErrorResponse   "Cart item not found"
// @Failure      500  {object}  dto.ErrorResponse   "Internal server error"
// @Router       /cart/items/{itemId} [delete]
// @Security     BearerAuth
func (h *CartHandler) RemoveItem(c *gin.Context) {

	userIDStr := c.GetString("userID")
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "user not authenticated"})
		return
	}

	// Get item ID from URL parameter
	itemIDParam := c.Param("itemId")
	var itemID uint
	if id, err := strconv.ParseUint(itemIDParam, 10, 32); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid item ID"})
		return
	} else {
		itemID = uint(id)
	}

	// Remove item from cart
	cartResponse, err := h.cartService.RemoveItem(userIDStr, itemID)
	if err != nil {
		if err.Error() == "cart item not found" || err.Error() == "cart item not found in user's cart" {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
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
// @Failure      401  {object}  dto.ErrorResponse   "Unauthorized"
// @Failure      404  {object}  dto.ErrorResponse   "Cart not found"
// @Failure      500  {object}  dto.ErrorResponse   "Internal server error"
// @Router       /cart [delete]
// @Security     BearerAuth
func (h *CartHandler) ClearCart(c *gin.Context) {

	userIDStr := c.GetString("userID")
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "user not authenticated"})
		return
	}

	// Clear all items from cart
	cartResponse, err := h.cartService.ClearCart(userIDStr)
	if err != nil {
		if err.Error() == "cart not found" {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to clear cart"})
		return
	}

	c.JSON(http.StatusOK, cartResponse)
}
