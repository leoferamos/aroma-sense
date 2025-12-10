package bootstrap

import (
	"github.com/leoferamos/aroma-sense/internal/handler"
	admin "github.com/leoferamos/aroma-sense/internal/handler/admin"
	product "github.com/leoferamos/aroma-sense/internal/handler/product"
	shipping "github.com/leoferamos/aroma-sense/internal/handler/shipping"
	"github.com/leoferamos/aroma-sense/internal/rate"
)

// initializeHandlers creates all handler instances
func initializeHandlers(services *services, rateLimiter rate.RateLimiter) *AppHandlers {

	return &AppHandlers{
		UserHandler:              handler.NewUserHandler(services.auth, services.userProfile, services.lgpd, services.chat),
		AdminUserHandler:         admin.NewAdminUserHandler(services.adminUser),
		ProductHandler:           product.NewProductHandler(services.product, services.review, services.userProfile),
		CartHandler:              handler.NewCartHandler(services.cart),
		OrderHandler:             handler.NewOrderHandler(services.order),
		PasswordResetHandler:     handler.NewPasswordResetHandler(services.passwordReset, rateLimiter),
		ShippingHandler:          shipping.NewShippingHandler(services.shipping),
		ReviewHandler:            product.NewReviewHandler(services.review, services.reviewReport, services.userProfile, services.product, services.auditLog, rateLimiter),
		AIHandler:                handler.NewAIHandler(services.ai, rateLimiter),
		ChatHandler:              handler.NewChatHandler(services.chat, rateLimiter),
		AuditLogHandler:          handler.NewAuditLogHandler(services.auditLog),
		AdminContestationHandler: admin.NewAdminContestationHandler(services.userContestation),
		AdminReviewReportHandler: admin.NewAdminReviewReportHandler(services.reviewReport),
		PaymentHandler:           handler.NewPaymentHandler(services.payment),
	}
}
