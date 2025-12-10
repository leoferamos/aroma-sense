package shipping

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/dto"
	handlererrors "github.com/leoferamos/aroma-sense/internal/handler/errors"
	shippingservice "github.com/leoferamos/aroma-sense/internal/service/shipping"
)

// ShippingHandler handles shipping quotation endpoints.
type ShippingHandler struct {
	shippingService shippingservice.ShippingService
}

func NewShippingHandler(shippingService shippingservice.ShippingService) *ShippingHandler {
	return &ShippingHandler{shippingService: shippingService}
}

// GetShippingOptions returns shipping options for the authenticated user's cart.
// @Summary      List shipping options
// @Description  Returns a list of shipping options (carrier, service code, price, ETA) for the current cart.
// @Tags         shipping
// @Produce      json
// @Param        postal_code query string true "Destination postal code (CEP)"
// @Success      200   {array} dto.ShippingOption
// @Failure      400   {object} dto.ErrorResponse "Error code: invalid_request"
// @Failure      401   {object} dto.ErrorResponse "Error code: unauthenticated"
// @Failure      500   {object} dto.ErrorResponse "Error code: internal_error"
// @Router       /shipping/options [get]
// @Security     BearerAuth
func (h *ShippingHandler) GetShippingOptions(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthenticated"})
		return
	}
	postalCode := c.Query("postal_code")
	if postalCode == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request"})
		return
	}
	options, err := h.shippingService.CalculateOptions(c.Request.Context(), userID, postalCode)
	if err != nil {
		if status, code, ok := handlererrors.MapServiceError(err); ok {
			c.JSON(status, dto.ErrorResponse{Error: code})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error"})
		return
	}
	c.JSON(http.StatusOK, options)
}
