package router

import (
	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/auth"
	paymenthandler "github.com/leoferamos/aroma-sense/internal/handler/payment"
)

// PaymentRoutes defines payment-related routes.
func PaymentRoutes(r *gin.Engine, handler *paymenthandler.PaymentHandler) {

	payments := r.Group("/payments")
	payments.Use(auth.JWTAuthMiddleware())
	{
		payments.POST("/intent", handler.CreateIntent)
	}

	// Payment webhook
	r.POST("/payments/webhook", handler.HandleWebhook)
}
