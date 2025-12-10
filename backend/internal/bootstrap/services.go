package bootstrap

import (
	"os"

	"github.com/leoferamos/aroma-sense/internal/notification"
	serviceadmin "github.com/leoferamos/aroma-sense/internal/service/admin"
	authservice "github.com/leoferamos/aroma-sense/internal/service/auth"
	cartservice "github.com/leoferamos/aroma-sense/internal/service/cart"
	chatservice "github.com/leoferamos/aroma-sense/internal/service/chat"
	lgpdservice "github.com/leoferamos/aroma-sense/internal/service/lgpd"
	logservice "github.com/leoferamos/aroma-sense/internal/service/log"
	orderservice "github.com/leoferamos/aroma-sense/internal/service/order"
	paymentservice "github.com/leoferamos/aroma-sense/internal/service/payment"
	productservice "github.com/leoferamos/aroma-sense/internal/service/product"
	reviewservice "github.com/leoferamos/aroma-sense/internal/service/review"
	shippingservice "github.com/leoferamos/aroma-sense/internal/service/shipping"
	userservice "github.com/leoferamos/aroma-sense/internal/service/user"
	"github.com/leoferamos/aroma-sense/internal/storage"
)

// services holds all service instances
type services struct {
	adminUser        serviceadmin.AdminUserService
	auth             authservice.AuthService
	userProfile      userservice.UserProfileService
	lgpd             lgpdservice.LgpdService
	product          productservice.ProductService
	cart             cartservice.CartService
	order            orderservice.OrderService
	payment          paymentservice.PaymentService
	passwordReset    authservice.PasswordResetService
	review           reviewservice.ReviewService
	reviewReport     reviewservice.ReviewReportService
	ai               *chatservice.AIService
	chat             *chatservice.ChatService
	shipping         shippingservice.ShippingService
	auditLog         logservice.AuditLogService
	userContestation userservice.UserContestationService
}

// initializeServices creates all service instances with proper dependencies
func initializeServices(repos *repositories, integrations *integrations, storageClient storage.ImageStorage) *services {
	frontend := os.Getenv("FRONTEND_URL")
	notifier := notification.NewNotifier(integrations.email, frontend)

	// Optional integrations
	if integrations.shipping != nil && integrations.shipping.provider != nil {
		integrations.shipping.service = shippingservice.NewShippingService(repos.cart, integrations.shipping.provider, integrations.shipping.originCEP)
	}

	// Core services in dependency order
	auditLogService := logservice.NewAuditLogService(repos.auditLog)
	aiService := chatservice.NewAIService(repos.product)
	productService := productservice.NewProductService(repos.product, storageClient, integrations.ai.embProvider)
	cartService := cartservice.NewCartService(repos.cart, productService)
	adminUserService := serviceadmin.NewAdminUserService(repos.user, auditLogService, notifier)
	userContestationService := userservice.NewUserContestationService(repos.userContestation, repos.user, adminUserService)
	lgpdService := lgpdservice.NewLgpdService(repos.user, repos.userContestation, auditLogService, notifier)
	reviewService := reviewservice.NewReviewService(repos.review, repos.order, repos.product)
	reviewReportService := reviewservice.NewReviewReportService(repos.reviewReport, repos.review, repos.user, adminUserService)
	chatService := chatservice.NewChatService(repos.product, integrations.ai.llmProvider, integrations.ai.embProvider, aiService)
	orderService := orderservice.NewOrderService(repos.order, repos.cart, repos.product, integrations.shipping.service)
	passwordResetService := authservice.NewPasswordResetService(repos.resetToken, repos.user, notifier)
	userProfileService := userservice.NewUserProfileService(repos.user, auditLogService)
	authService := authservice.NewAuthService(repos.user, cartService, auditLogService)

	var paymentSvc paymentservice.PaymentService
	if integrations.payment != nil && integrations.payment.provider != nil {
		paymentSvc = paymentservice.NewPaymentService(repos.cart, repos.product, repos.order, repos.payment, integrations.shipping.service, integrations.payment.provider)
	}

	return &services{
		adminUser:        adminUserService,
		auth:             authService,
		userProfile:      userProfileService,
		lgpd:             lgpdService,
		product:          productService,
		cart:             cartService,
		order:            orderService,
		payment:          paymentSvc,
		passwordReset:    passwordResetService,
		review:           reviewService,
		reviewReport:     reviewReportService,
		ai:               aiService,
		chat:             chatService,
		shipping:         integrations.shipping.service,
		auditLog:         auditLogService,
		userContestation: userContestationService,
	}
}
