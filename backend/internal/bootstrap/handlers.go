package bootstrap

import (
	"github.com/leoferamos/aroma-sense/internal/handler"
	"github.com/leoferamos/aroma-sense/internal/rate"
)

// initializeHandlers creates all handler instances
func initializeHandlers(services *services, rateLimiter rate.RateLimiter) *AppHandlers {

	return &AppHandlers{
		UserHandler:              handler.NewUserHandler(services.auth, services.userProfile, services.lgpd),
		AdminUserHandler:         handler.NewAdminUserHandler(services.adminUser),
		ProductHandler:           handler.NewProductHandler(services.product, services.review, services.userProfile),
		CartHandler:              handler.NewCartHandler(services.cart),
		OrderHandler:             handler.NewOrderHandler(services.order),
		PasswordResetHandler:     handler.NewPasswordResetHandler(services.passwordReset, rateLimiter),
		ShippingHandler:          handler.NewShippingHandler(services.shipping),
		ReviewHandler:            handler.NewReviewHandler(services.review, services.reviewReport, services.userProfile, services.product, services.auditLog, rateLimiter),
		AIHandler:                handler.NewAIHandler(services.ai, rateLimiter),
		ChatHandler:              handler.NewChatHandler(services.chat, rateLimiter),
		AuditLogHandler:          handler.NewAuditLogHandler(services.auditLog),
		AdminContestationHandler: handler.NewAdminContestationHandler(services.userContestation),
		AdminReviewReportHandler: handler.NewAdminReviewReportHandler(services.reviewReport),
		PaymentHandler:           handler.NewPaymentHandler(services.payment),
	}
}
