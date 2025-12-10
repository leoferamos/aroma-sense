package router

import (
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/bootstrap"
	"github.com/leoferamos/aroma-sense/internal/middleware"
)

// SetupRouter initializes the Gin router with all routes
func SetupRouter(handlers *bootstrap.AppHandlers) *gin.Engine {
	r := gin.Default()

	// CORS setup: read allowed origins from env
	allowedOrigins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	// Health check
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Provide user profile service to auth middleware
	middleware.SetUserProfileService(handlers.UserHandler.UserProfile())

	// Register domain routes
	UserRoutes(r, handlers.UserHandler, handlers.PasswordResetHandler)
	AdminRoutes(r, handlers.AdminUserHandler, handlers.ProductHandler, handlers.OrderHandler, handlers.AuditLogHandler, handlers.AdminContestationHandler, handlers.AdminReviewReportHandler)
	ProductRoutes(r, handlers.ProductHandler, handlers.ReviewHandler)
	CartRoutes(r, handlers.CartHandler)
	OrderRoutes(r, handlers.OrderHandler)
	ShippingRoutes(r, handlers.ShippingHandler)
	AIRoutes(r, handlers.AIHandler, handlers.ChatHandler)
	PaymentRoutes(r, handlers.PaymentHandler)

	return r
}
