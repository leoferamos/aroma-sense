package router

import (
	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/auth"
	"github.com/leoferamos/aroma-sense/internal/handler"
)

// ShippingRoutes sets up the shipping-related routes
func ShippingRoutes(r *gin.Engine, shippingHandler *handler.ShippingHandler) {
	grp := r.Group("/shipping")
	grp.Use(auth.JWTAuthMiddleware())
	{
		grp.GET("/options", shippingHandler.GetShippingOptions)
	}
}
