package router

import (
	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/auth"
	"github.com/leoferamos/aroma-sense/internal/handler"
	"github.com/leoferamos/aroma-sense/internal/middleware"
)

// ProductRoutes sets up the product-related routes
func ProductRoutes(r *gin.Engine, productHandler *handler.ProductHandler, reviewHandler *handler.ReviewHandler) {
	productGroup := r.Group("/products")
	productGroup.Use(auth.OptionalAuthMiddleware(), middleware.AccountStatusMiddleware())
	{
		productGroup.GET("", productHandler.GetLatestProducts)
		productGroup.GET("/:id/reviews/summary", reviewHandler.GetSummary)
		productGroup.POST("/:id/reviews", auth.JWTAuthMiddleware(), reviewHandler.CreateReview)
		productGroup.GET("/:id/reviews", reviewHandler.ListReviews)
		productGroup.GET("/:id", productHandler.GetProduct)
	}

	// Clean URLs for users
	productSlugGroup := r.Group("/product")
	productSlugGroup.Use(auth.OptionalAuthMiddleware(), middleware.AccountStatusMiddleware())
	{
		productSlugGroup.GET("/:slug", productHandler.GetProductBySlug)
		productSlugGroup.GET("/:slug/reviews/summary", reviewHandler.GetSummary)
		productSlugGroup.POST("/:slug/reviews", auth.JWTAuthMiddleware(), reviewHandler.CreateReview)
		productSlugGroup.GET("/:slug/reviews", reviewHandler.ListReviews)
	}
}
