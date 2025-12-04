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
	// Initialize shipping service if provider is available
	if integrations.shipping != nil && integrations.shipping.provider != nil {
		integrations.shipping.service = service.NewShippingService(
			repos.cart,
			integrations.shipping.provider,
			integrations.shipping.originCEP,
		)
	}

	// Initialize audit log service first
	auditLogService := service.NewAuditLogService(repos.auditLog)

	// Initialize services in dependency order
	cartService := service.NewCartService(repos.cart, nil)
	// create notifier to encapsulate email sending
	frontend := os.Getenv("FRONTEND_URL")
	notifier := notification.NewNotifier(integrations.email, frontend)
	adminUserService := service.NewAdminUserService(repos.user, auditLogService, notifier)
	authService := service.NewAuthService(repos.user, cartService, auditLogService)
	userProfileService := service.NewUserProfileService(repos.user, auditLogService)
	lgpdService := service.NewLgpdService(repos.user, repos.userContestation, auditLogService, notifier)
	orderService := service.NewOrderService(repos.order, repos.cart, repos.product, integrations.shipping.service)
	passwordResetService := service.NewPasswordResetService(repos.resetToken, repos.user, notifier)
	reviewService := service.NewReviewService(repos.review, repos.order, repos.product)
	aiService := service.NewAIService(repos.product)

	productService := service.NewProductService(repos.product, storageClient, integrations.ai.embProvider)
	cartService = service.NewCartService(repos.cart, productService)
	chatService := service.NewChatService(repos.product, integrations.ai.llmProvider, integrations.ai.embProvider, aiService)

	userContestationService := service.NewUserContestationService(repos.userContestation)
	return &services{
		adminUser:        adminUserService,
		auth:             authService,
		userProfile:      userProfileService,
		lgpd:             lgpdService,
		product:          productService,
		cart:             cartService,
		order:            orderService,
		passwordReset:    passwordResetService,
		review:           reviewService,
		ai:               aiService,
		chat:             chatService,
		shipping:         integrations.shipping.service,
		auditLog:         auditLogService,
		userContestation: userContestationService,
	}
}
