package bootstrap

import (
	"github.com/leoferamos/aroma-sense/internal/email"
	"github.com/leoferamos/aroma-sense/internal/handler"
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
	emailService, err := email.NewSMTPEmailService(smtpConfig)
	if err != nil {
		emailService = nil
	}

	// Initialize services in dependency order
	productService := service.NewProductService(productRepo, storageClient)
	cartService := service.NewCartService(cartRepo, productService)
	userService := service.NewUserService(userRepo, cartService)
	orderService := service.NewOrderService(orderRepo, cartRepo, productRepo)
	passwordResetService := service.NewPasswordResetService(resetTokenRepo, userRepo, emailService)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)
	productHandler := handler.NewProductHandler(productService)
	cartHandler := handler.NewCartHandler(cartService)
	orderHandler := handler.NewOrderHandler(orderService)
	passwordResetHandler := handler.NewPasswordResetHandler(passwordResetService, rateLimiter)

	return &AppHandlers{
		UserHandler:          userHandler,
		ProductHandler:       productHandler,
		CartHandler:          cartHandler,
		OrderHandler:         orderHandler,
		PasswordResetHandler: passwordResetHandler,
	}
}
