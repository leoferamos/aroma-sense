package bootstrap

import (
	"github.com/leoferamos/aroma-sense/internal/handler"
	"github.com/leoferamos/aroma-sense/internal/rate"
)

// initializeHandlers creates all handler instances
func initializeHandlers(services *services, rateLimiter rate.RateLimiter) *AppHandlers {
	return &AppHandlers{
		UserHandler:          handler.NewUserHandler(services.user),
		ProductHandler:       handler.NewProductHandler(services.product, services.review, services.user),
		CartHandler:          handler.NewCartHandler(services.cart),
		OrderHandler:         handler.NewOrderHandler(services.order),
		PasswordResetHandler: handler.NewPasswordResetHandler(services.passwordReset, rateLimiter),
		ShippingHandler:      handler.NewShippingHandler(services.shipping),
		ReviewHandler:        handler.NewReviewHandler(services.review, services.user),
		AIHandler:            handler.NewAIHandler(services.ai, rateLimiter),
		ChatHandler:          handler.NewChatHandler(services.chat, rateLimiter),
	}
}
