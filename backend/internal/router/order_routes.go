package router

import (
	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/auth"
	orderhandler "github.com/leoferamos/aroma-sense/internal/handler/order"
	"github.com/leoferamos/aroma-sense/internal/middleware"
)

// OrderRoutes sets up the order-related routes
func OrderRoutes(r *gin.Engine, orderHandler *orderhandler.OrderHandler) {
	orderGroup := r.Group("/orders")
	orderGroup.Use(auth.JWTAuthMiddleware(), middleware.AccountStatusMiddleware())
	{
		orderGroup.GET("", orderHandler.ListUserOrders)
		orderGroup.POST("", orderHandler.CreateOrderFromCart)
	}
}
