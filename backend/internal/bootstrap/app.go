package bootstrap

import (
	"github.com/leoferamos/aroma-sense/internal/email"
	"github.com/leoferamos/aroma-sense/internal/handler"
	shippingprovider "github.com/leoferamos/aroma-sense/internal/integrations/shipping"
	"github.com/leoferamos/aroma-sense/internal/rate"
	"github.com/leoferamos/aroma-sense/internal/repository"
	"github.com/leoferamos/aroma-sense/internal/service"
	"github.com/leoferamos/aroma-sense/internal/storage"
	"gorm.io/gorm"
)

// AppHandlers contains all initialized handlers
type AppHandlers struct {
	UserHandler          *handler.UserHandler
	ProductHandler       *handler.ProductHandler
	CartHandler          *handler.CartHandler
	OrderHandler         *handler.OrderHandler
	PasswordResetHandler *handler.PasswordResetHandler
	ShippingHandler      *handler.ShippingHandler
	ReviewHandler        *handler.ReviewHandler
}

// InitializeApp initializes all modules with proper dependency injection
func InitializeApp(db *gorm.DB, storageClient storage.ImageStorage) *AppHandlers {
	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	cartRepo := repository.NewCartRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	resetTokenRepo := repository.NewResetTokenRepository(db)

	// Initialize rate limiter
	rateLimiter := rate.NewInMemory()

	// Initialize email service
	smtpConfig := email.LoadSMTPConfigFromEnv()
	if err := smtpConfig.Validate(); err != nil {
		panic("SMTP configuration error: " + err.Error())
	}
	emailService, err := email.NewSMTPEmailService(smtpConfig)
	if err != nil {
		panic("Failed to initialize email service: " + err.Error())
	}

	// Initialize shipping provider and service
	var shippingProvider service.ShippingProvider
	var shippingService service.ShippingService
	if cfg, err := shippingprovider.LoadShippingConfigFromEnv(); err == nil {
		if cli, err := shippingprovider.NewClient(cfg); err == nil {
			provider := shippingprovider.NewProvider(cli).
				WithQuotesPath(cfg.QuotesPath).
				WithStaticAuth(cfg.StaticToken, cfg.UserAgent).
				WithServices(cfg.Services)
			shippingProvider = provider
			shippingService = service.NewShippingService(cartRepo, shippingProvider, cfg.OriginCEP)
		}
	}

	// Initialize services in dependency order
	productService := service.NewProductService(productRepo, storageClient)
	cartService := service.NewCartService(cartRepo, productService)
	userService := service.NewUserService(userRepo, cartService)
	orderService := service.NewOrderService(orderRepo, cartRepo, productRepo, shippingService)
	passwordResetService := service.NewPasswordResetService(resetTokenRepo, userRepo, emailService)
	reviewRepo := repository.NewReviewRepository(db)
	reviewService := service.NewReviewService(reviewRepo, orderRepo, productRepo)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)
	productHandler := handler.NewProductHandler(productService, reviewService, userRepo)
	reviewHandler := handler.NewReviewHandler(reviewService, userRepo)
	cartHandler := handler.NewCartHandler(cartService)
	orderHandler := handler.NewOrderHandler(orderService)
	shippingHandler := handler.NewShippingHandler(shippingService)
	passwordResetHandler := handler.NewPasswordResetHandler(passwordResetService, rateLimiter)

	return &AppHandlers{
		UserHandler:          userHandler,
		ProductHandler:       productHandler,
		CartHandler:          cartHandler,
		OrderHandler:         orderHandler,
		PasswordResetHandler: passwordResetHandler,
		ShippingHandler:      shippingHandler,
		ReviewHandler:        reviewHandler,
	}
}
