package bootstrap

import (
	"github.com/leoferamos/aroma-sense/internal/handler"
	"github.com/leoferamos/aroma-sense/internal/repository"
	"github.com/leoferamos/aroma-sense/internal/service"
	"github.com/leoferamos/aroma-sense/internal/storage"
	"gorm.io/gorm"
)

// AppHandlers contains all initialized handlers
type AppHandlers struct {
	UserHandler    *handler.UserHandler
	ProductHandler *handler.ProductHandler
	CartHandler    *handler.CartHandler
}

// InitializeApp initializes all modules with proper dependency injection
func InitializeApp(db *gorm.DB, storageClient storage.ImageStorage) *AppHandlers {
	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	cartRepo := repository.NewCartRepository(db)

	// Initialize services in dependency order
	productService := service.NewProductService(productRepo, storageClient)
	cartService := service.NewCartService(cartRepo, productService)
	userService := service.NewUserService(userRepo, cartService)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)
	productHandler := handler.NewProductHandler(productService)
	cartHandler := handler.NewCartHandler(cartService)

	return &AppHandlers{
		UserHandler:    userHandler,
		ProductHandler: productHandler,
		CartHandler:    cartHandler,
	}
}
