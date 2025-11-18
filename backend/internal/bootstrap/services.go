package bootstrap

import (
	"github.com/leoferamos/aroma-sense/internal/service"
	"github.com/leoferamos/aroma-sense/internal/storage"
)

// services holds all service instances
type services struct {
	user          service.UserService
	product       service.ProductService
	cart          service.CartService
	order         service.OrderService
	passwordReset service.PasswordResetService
	review        service.ReviewService
	ai            *service.AIService
	chat          *service.ChatService
	shipping      service.ShippingService
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

	// Initialize services in dependency order
	cartService := service.NewCartService(repos.cart, nil)
	userService := service.NewUserService(repos.user, cartService)
	orderService := service.NewOrderService(repos.order, repos.cart, repos.product, integrations.shipping.service)
	passwordResetService := service.NewPasswordResetService(repos.resetToken, repos.user, integrations.email)
	reviewService := service.NewReviewService(repos.review, repos.order, repos.product)
	aiService := service.NewAIService(repos.product)

	productService := service.NewProductService(repos.product, storageClient, integrations.ai.embProvider)
	cartService = service.NewCartService(repos.cart, productService)
	chatService := service.NewChatService(repos.product, integrations.ai.llmProvider, integrations.ai.embProvider)

	return &services{
		user:          userService,
		product:       productService,
		cart:          cartService,
		order:         orderService,
		passwordReset: passwordResetService,
		review:        reviewService,
		ai:            aiService,
		chat:          chatService,
		shipping:      integrations.shipping.service,
	}
}
