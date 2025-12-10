package bootstrap

import (
	admin "github.com/leoferamos/aroma-sense/internal/handler/admin"
	auth "github.com/leoferamos/aroma-sense/internal/handler/auth"
	carthandler "github.com/leoferamos/aroma-sense/internal/handler/cart"
	chathandler "github.com/leoferamos/aroma-sense/internal/handler/chat"
	loghandler "github.com/leoferamos/aroma-sense/internal/handler/log"
	orderhandler "github.com/leoferamos/aroma-sense/internal/handler/order"
	paymenthandler "github.com/leoferamos/aroma-sense/internal/handler/payment"
	product "github.com/leoferamos/aroma-sense/internal/handler/product"
	reviewhandler "github.com/leoferamos/aroma-sense/internal/handler/review"
	shipping "github.com/leoferamos/aroma-sense/internal/handler/shipping"
	userhandler "github.com/leoferamos/aroma-sense/internal/handler/user"
	"github.com/leoferamos/aroma-sense/internal/rate"
)

// initializeHandlers creates all handler instances
func initializeHandlers(services *services, rateLimiter rate.RateLimiter) *AppHandlers {

	return &AppHandlers{
		UserHandler:              userhandler.NewUserHandler(services.auth, services.userProfile, services.lgpd, services.chat),
		AdminUserHandler:         admin.NewAdminUserHandler(services.adminUser),
		ProductHandler:           product.NewProductHandler(services.product, services.review, services.userProfile),
		CartHandler:              carthandler.NewCartHandler(services.cart),
		OrderHandler:             orderhandler.NewOrderHandler(services.order),
		PasswordResetHandler:     auth.NewPasswordResetHandler(services.passwordReset, rateLimiter),
		ShippingHandler:          shipping.NewShippingHandler(services.shipping),
		ReviewHandler:            reviewhandler.NewReviewHandler(services.review, services.reviewReport, services.userProfile, services.product, services.auditLog, rateLimiter),
		AIHandler:                chathandler.NewAIHandler(services.ai, rateLimiter),
		ChatHandler:              chathandler.NewChatHandler(services.chat, rateLimiter),
		AuditLogHandler:          loghandler.NewAuditLogHandler(services.auditLog),
		AdminContestationHandler: admin.NewAdminContestationHandler(services.userContestation),
		AdminReviewReportHandler: admin.NewAdminReviewReportHandler(services.reviewReport),
		PaymentHandler:           paymenthandler.NewPaymentHandler(services.payment),
	}
}
