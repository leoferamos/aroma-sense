package bootstrap

import (
	"github.com/leoferamos/aroma-sense/internal/handler"
	"github.com/leoferamos/aroma-sense/internal/rate"
	"github.com/leoferamos/aroma-sense/internal/storage"
	"gorm.io/gorm"
)

// AppHandlers contains all initialized handlers
type AppHandlers struct {
	UserHandler          *handler.UserHandler
	AdminUserHandler     *handler.AdminUserHandler
	ProductHandler       *handler.ProductHandler
	CartHandler          *handler.CartHandler
	OrderHandler         *handler.OrderHandler
	PasswordResetHandler *handler.PasswordResetHandler
	ShippingHandler      *handler.ShippingHandler
	ReviewHandler        *handler.ReviewHandler
	AIHandler            *handler.AIHandler
	ChatHandler          *handler.ChatHandler
}

// InitializeApp initializes all modules with proper dependency injection
func InitializeApp(db *gorm.DB, storageClient storage.ImageStorage) *AppHandlers {
	// Initialize core infrastructure
	repositories := initializeRepositories(db)
	rateLimiter := rate.NewInMemory()

	// Initialize external integrations
	integrations := initializeIntegrations()

	// Initialize services in dependency order
	services := initializeServices(repositories, integrations, storageClient)

	// Initialize handlers
	handlers := initializeHandlers(services, rateLimiter)

	return handlers
}
