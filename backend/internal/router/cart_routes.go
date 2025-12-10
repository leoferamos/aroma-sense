package router

import (
	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/auth"
	carthandler "github.com/leoferamos/aroma-sense/internal/handler/cart"
	"github.com/leoferamos/aroma-sense/internal/middleware"
)

// CartRoutes sets up the cart-related routes
func CartRoutes(r *gin.Engine, handler *carthandler.CartHandler) {
	cartGroup := r.Group("/cart")
	cartGroup.Use(auth.JWTAuthMiddleware(), middleware.AccountStatusMiddleware())
	{
		cartGroup.GET("", handler.GetCart)
		cartGroup.POST("", handler.AddItem)
		cartGroup.DELETE("", handler.ClearCart)
		cartGroup.PATCH("/items/:productSlug", handler.UpdateItemQuantity)
		cartGroup.DELETE("/items/:productSlug", handler.RemoveItem)
	}
}
