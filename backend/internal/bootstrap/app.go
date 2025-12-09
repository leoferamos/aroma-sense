package bootstrap

import (
	"github.com/leoferamos/aroma-sense/internal/handler"
	"github.com/leoferamos/aroma-sense/internal/rate"
	"github.com/leoferamos/aroma-sense/internal/repository"
	"github.com/leoferamos/aroma-sense/internal/service"
	"github.com/leoferamos/aroma-sense/internal/storage"
	"gorm.io/gorm"
)

// AppHandlers contains all initialized handlers
type AppHandlers struct {
	UserHandler              *handler.UserHandler
	AdminUserHandler         *handler.AdminUserHandler
	ProductHandler           *handler.ProductHandler
	CartHandler              *handler.CartHandler
	OrderHandler             *handler.OrderHandler
	PasswordResetHandler     *handler.PasswordResetHandler
	ShippingHandler          *handler.ShippingHandler
	ReviewHandler            *handler.ReviewHandler
	AIHandler                *handler.AIHandler
	ChatHandler              *handler.ChatHandler
	AuditLogHandler          *handler.AuditLogHandler
	AdminContestationHandler *handler.AdminContestationHandler
	PaymentHandler           *handler.PaymentHandler
}

// AppServices contains service instances needed for jobs
type AppServices struct {
	AdminUserService service.AdminUserService
	AuditLogService  service.AuditLogService
	LgpdService      service.LgpdService
}

// AppRepos contains repository instances needed for jobs
type AppRepos struct {
	UserRepo repository.UserRepository
}

// AppComponents contains all initialized application components
type AppComponents struct {
	Handlers *AppHandlers
	Services *AppServices
	Repos    *AppRepos
}

// InitializeApp initializes all modules with proper dependency injection
func InitializeApp(db *gorm.DB, storageClient storage.ImageStorage) *AppComponents {
	// Initialize core infrastructure
	repositories := initializeRepositories(db)
	rateLimiter := rate.NewInMemory()

	// Initialize external integrations
	integrations := initializeIntegrations()

	// Initialize services in dependency order
	services := initializeServices(repositories, integrations, storageClient)

	// Initialize handlers
	handlers := initializeHandlers(services, rateLimiter)

	appServices := &AppServices{
		AdminUserService: services.adminUser,
		AuditLogService:  services.auditLog,
		LgpdService:      services.lgpd,
	}

	appRepos := &AppRepos{
		UserRepo: repositories.user,
	}

	return &AppComponents{
		Handlers: handlers,
		Services: appServices,
		Repos:    appRepos,
	}
}
