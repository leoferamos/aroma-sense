package router

import (
	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/handler"
)

func OrderRoutes(r *gin.Engine, orderHandler *handler.OrderHandler) {
	orderGroup := r.Group("/orders")
	{
		orderGroup.POST("", orderHandler.CreateOrderFromCart)
	}
}
