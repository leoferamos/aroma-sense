package router

import (
	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/auth"
	"github.com/leoferamos/aroma-sense/internal/handler"
)

// CartRoutes sets up the cart-related routes
func CartRoutes(r *gin.Engine, handler *handler.CartHandler) {
	cartGroup := r.Group("/cart")
	cartGroup.Use(auth.JWTAuthMiddleware())
	{
		cartGroup.GET("", handler.GetCart)
		cartGroup.POST("", handler.AddItem)
		cartGroup.PATCH("/items/:itemId", handler.UpdateItemQuantity)
		cartGroup.DELETE("/items/:itemId", handler.RemoveItem)
	}
}
