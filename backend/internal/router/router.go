package router

import (
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/handler"
)

// SetupRouter initializes the Gin router with all routes
func SetupRouter(userHandler *handler.UserHandler, productHandler *handler.ProductHandler, cartHandler *handler.CartHandler) *gin.Engine {
	r := gin.Default()

	// CORS setup: read allowed origins from env
	allowedOrigins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	// Health check
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Register domain routes
	UserRoutes(r, userHandler)
	AdminRoutes(r, userHandler, productHandler)
	ProductRoutes(r, productHandler)
	CartRoutes(r, cartHandler)

	return r
}
