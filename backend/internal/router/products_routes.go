package router

import (
	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/auth"
	"github.com/leoferamos/aroma-sense/internal/handler"
	"github.com/leoferamos/aroma-sense/internal/middleware"
)

// ProductRoutes sets up the product-related routes
func ProductRoutes(r *gin.Engine, productHandler *handler.ProductHandler, reviewHandler *handler.ReviewHandler) {
	// Public routes - always use slug
	productGroup := r.Group("/products")
	productGroup.Use(auth.OptionalAuthMiddleware(), middleware.AccountStatusMiddleware())
	{
		productGroup.GET("", productHandler.GetLatestProducts)
		productGroup.GET("/:slug/reviews/summary", reviewHandler.GetSummary)
		productGroup.POST("/:slug/reviews", auth.JWTAuthMiddleware(), reviewHandler.CreateReview)
		productGroup.GET("/:slug/reviews", reviewHandler.ListReviews)
		productGroup.GET("/:slug", productHandler.GetProduct)
	}
}
