package router

import (
	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/handler"
)

// ProductRoutes sets up the product-related routes
func ProductRoutes(r *gin.Engine, handler *handler.ProductHandler) {
	productGroup := r.Group("/products")
	{
		productGroup.GET("/:id", handler.GetProduct)
		productGroup.GET("/", handler.GetProducts)
	}
}
