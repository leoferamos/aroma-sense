package router

import (
	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/auth"
	"github.com/leoferamos/aroma-sense/internal/handler"
	"github.com/leoferamos/aroma-sense/internal/middleware"
)

// CartRoutes sets up the cart-related routes
func CartRoutes(r *gin.Engine, handler *handler.CartHandler) {
	cartGroup := r.Group("/cart")
	cartGroup.Use(auth.JWTAuthMiddleware(), middleware.AccountStatusMiddleware())
	{
		cartGroup.GET("", handler.GetCart)
		cartGroup.POST("", handler.AddItem)
		cartGroup.DELETE("", handler.ClearCart)
		cartGroup.PATCH("/items/:itemId", handler.UpdateItemQuantity)
		cartGroup.DELETE("/items/:itemId", handler.RemoveItem)
	}
}
