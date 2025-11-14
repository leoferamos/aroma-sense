package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/dto"
	"github.com/leoferamos/aroma-sense/internal/service"
)

// ShippingHandler handles shipping quotation endpoints.
type ShippingHandler struct {
	shippingService service.ShippingService
}

func NewShippingHandler(shippingService service.ShippingService) *ShippingHandler {
	return &ShippingHandler{shippingService: shippingService}
}

// GetShippingOptions returns shipping options for the authenticated user's cart.
// @Summary      List shipping options
// @Description  Returns a list of shipping options (carrier, service code, price, ETA) for the current cart.
// @Tags         shipping
// @Produce      json
// @Param        postal_code query string true "Destination postal code (CEP)"
// @Success      200   {array} dto.ShippingOption
// @Failure      400   {object} dto.ErrorResponse
// @Failure      401   {object} dto.ErrorResponse
// @Failure      500   {object} dto.ErrorResponse
// @Router       /shipping/options [get]
// @Security     BearerAuth
func (h *ShippingHandler) GetShippingOptions(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}
	postalCode := c.Query("postal_code")
	if postalCode == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "postal_code required"})
		return
	}
	options, err := h.shippingService.CalculateOptions(c.Request.Context(), userID, postalCode)
	if err != nil {
		if status, msg, ok := mapServiceError(err); ok {
			c.JSON(status, dto.ErrorResponse{Error: msg})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to quote shipping"})
		return
	}
	c.JSON(http.StatusOK, options)
}
