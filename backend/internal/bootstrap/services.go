package bootstrap

import (
	"os"

	"github.com/leoferamos/aroma-sense/internal/notification"
	"github.com/leoferamos/aroma-sense/internal/service"
	"github.com/leoferamos/aroma-sense/internal/storage"
)

// services holds all service instances
type services struct {
	adminUser        service.AdminUserService
	auth             service.AuthService
	userProfile      service.UserProfileService
	lgpd             service.LgpdService
	product          service.ProductService
	cart             service.CartService
	order            service.OrderService
	payment          service.PaymentService
	passwordReset    service.PasswordResetService
	review           service.ReviewService
	ai               *service.AIService
	chat             *service.ChatService
	shipping         service.ShippingService
	auditLog         service.AuditLogService
	userContestation service.UserContestationService
}

// initializeServices creates all service instances with proper dependencies
func initializeServices(repos *repositories, integrations *integrations, storageClient storage.ImageStorage) *services {
	frontend := os.Getenv("FRONTEND_URL")
	notifier := notification.NewNotifier(integrations.email, frontend)

	// Optional integrations
	if integrations.shipping != nil && integrations.shipping.provider != nil {
		integrations.shipping.service = service.NewShippingService(repos.cart, integrations.shipping.provider, integrations.shipping.originCEP)
	}

	// Core services in dependency order
	auditLogService := service.NewAuditLogService(repos.auditLog)
	aiService := service.NewAIService(repos.product)
	productService := service.NewProductService(repos.product, storageClient, integrations.ai.embProvider)
	cartService := service.NewCartService(repos.cart, productService)
	adminUserService := service.NewAdminUserService(repos.user, auditLogService, notifier)
	userContestationService := service.NewUserContestationService(repos.userContestation, repos.user, adminUserService)
	lgpdService := service.NewLgpdService(repos.user, repos.userContestation, auditLogService, notifier)
	reviewService := service.NewReviewService(repos.review, repos.order, repos.product)
	chatService := service.NewChatService(repos.product, integrations.ai.llmProvider, integrations.ai.embProvider, aiService)
	orderService := service.NewOrderService(repos.order, repos.cart, repos.product, integrations.shipping.service)
	passwordResetService := service.NewPasswordResetService(repos.resetToken, repos.user, notifier)
	userProfileService := service.NewUserProfileService(repos.user, auditLogService)
	authService := service.NewAuthService(repos.user, cartService, auditLogService)

	var paymentSvc service.PaymentService
	if integrations.payment != nil && integrations.payment.provider != nil {
		paymentSvc = service.NewPaymentService(repos.cart, repos.product, repos.order, repos.payment, integrations.shipping.service, integrations.payment.provider)
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
		ai:               aiService,
		chat:             chatService,
		shipping:         integrations.shipping.service,
		auditLog:         auditLogService,
		userContestation: userContestationService,
	}
}
