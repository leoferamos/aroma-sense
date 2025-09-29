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
	cart, err := h.cartService.GetCartByUserID(userIDStr)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "cart not found"})
		return
	}

	cartResponse := dto.CartResponse{
		ID:        cart.ID,
		UserID:    cart.UserID,
		Items:     []dto.CartItemResponse{},
		Total:     0.0,
		ItemCount: 0,
		CreatedAt: cart.CreatedAt,
		UpdatedAt: cart.UpdatedAt,
	}

	// Convert cart items and calculate totals
	for _, item := range cart.Items {
		itemTotal := item.Price * float64(item.Quantity)

		cartItemResponse := dto.CartItemResponse{
			ID:        item.ID,
			CartID:    item.CartID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
			Total:     itemTotal,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		}

		if item.Product != nil {
			cartItemResponse.Product = &dto.ProductResponse{
				ID:            item.Product.ID,
				Name:          item.Product.Name,
				Brand:         item.Product.Brand,
				Weight:        item.Product.Weight,
				Description:   item.Product.Description,
				Price:         item.Product.Price,
				ImageURL:      item.Product.ImageURL,
				Category:      item.Product.Category,
				Notes:         item.Product.Notes,
				StockQuantity: item.Product.StockQuantity,
				CreatedAt:     item.Product.CreatedAt,
				UpdatedAt:     item.Product.UpdatedAt,
			}
		}

		cartResponse.Items = append(cartResponse.Items, cartItemResponse)
		cartResponse.Total += itemTotal
		cartResponse.ItemCount += item.Quantity
	}

	c.JSON(http.StatusOK, cartResponse)
}
