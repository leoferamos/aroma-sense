package router

import (
	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/auth"
	"github.com/leoferamos/aroma-sense/internal/handler"
)

// ProductRoutes sets up the product-related routes
func ProductRoutes(r *gin.Engine, productHandler *handler.ProductHandler, reviewHandler *handler.ReviewHandler) {
	productGroup := r.Group("/products")
	productGroup.Use(auth.OptionalAuthMiddleware())
	{
		productGroup.GET("", productHandler.GetLatestProducts)
		productGroup.GET("/:id", productHandler.GetProductByID)
		productGroup.GET("/:id/reviews", reviewHandler.ListReviews)
		productGroup.GET("/:id/reviews/summary", reviewHandler.GetSummary)
		productGroup.POST("/:id/reviews", auth.JWTAuthMiddleware(), reviewHandler.CreateReview)
	}
}
