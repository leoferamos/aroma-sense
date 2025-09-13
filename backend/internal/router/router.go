package router

import (
	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/handler"
)

// SetupRouter initializes the Gin router with all routes
func SetupRouter(userHandler *handler.UserHandler, productHandler *handler.ProductHandler) *gin.Engine {
	r := gin.Default()

	// Health check
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Register domain routes
	UserRoutes(r, userHandler)
	AdminRoutes(r, userHandler, productHandler)
	ProductRoutes(r, productHandler)

	return r
}
