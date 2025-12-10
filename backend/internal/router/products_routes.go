package router

import (
	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/auth"
	product "github.com/leoferamos/aroma-sense/internal/handler/product"
	reviewhandler "github.com/leoferamos/aroma-sense/internal/handler/review"
	"github.com/leoferamos/aroma-sense/internal/middleware"
)

// ProductRoutes sets up the product-related routes
func ProductRoutes(r *gin.Engine, productHandler *product.ProductHandler, reviewHandler *reviewhandler.ReviewHandler) {
	// Public routes
	publicProductGroup := r.Group("/products")
	publicProductGroup.Use(auth.OptionalAuthMiddleware(), middleware.AccountStatusMiddleware())
	{
		// Product listing and details
		publicProductGroup.GET("", productHandler.GetLatestProducts)
		publicProductGroup.GET("/:slug", productHandler.GetProduct)

		// Public review operations
		publicProductGroup.GET("/:slug/reviews", reviewHandler.ListReviews)
		publicProductGroup.GET("/:slug/reviews/summary", reviewHandler.GetSummary)
	}

	// Authenticated routes
	authenticatedGroup := r.Group("")
	authenticatedGroup.Use(auth.JWTAuthMiddleware())
	{
		authenticatedGroup.POST("/products/:slug/reviews", reviewHandler.CreateReview)
		authenticatedGroup.DELETE("/reviews/:reviewID", reviewHandler.DeleteReview)
		authenticatedGroup.POST("/reviews/:reviewID/report", reviewHandler.ReportReview)
	}
}
