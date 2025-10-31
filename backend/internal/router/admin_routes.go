package router

import (
	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/auth"
	"github.com/leoferamos/aroma-sense/internal/handler"
)

// AdminRoutes sets up the admin-related routes
func AdminRoutes(r *gin.Engine, userHandler *handler.UserHandler,
	productHandler *handler.ProductHandler, orderHandler *handler.OrderHandler) {
	adminGroup := r.Group("/admin")
	adminGroup.Use(auth.JWTAuthMiddleware(), auth.AdminOnly())
	{
		adminGroup.POST("/products", productHandler.CreateProduct)
		adminGroup.GET("/products/:id", productHandler.GetProduct)
		adminGroup.PATCH("/products/:id", productHandler.UpdateProduct)
		adminGroup.DELETE("/products/:id", productHandler.DeleteProduct)
		adminGroup.GET("/orders", orderHandler.ListOrders)
	}
}
