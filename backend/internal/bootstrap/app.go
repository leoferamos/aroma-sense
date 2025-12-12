package bootstrap

import (
	admin "github.com/leoferamos/aroma-sense/internal/handler/admin"
	aihandler "github.com/leoferamos/aroma-sense/internal/handler/ai"
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
	"github.com/leoferamos/aroma-sense/internal/repository"
	serviceadmin "github.com/leoferamos/aroma-sense/internal/service/admin"
	servicelgpd "github.com/leoferamos/aroma-sense/internal/service/lgpd"
	servicelog "github.com/leoferamos/aroma-sense/internal/service/log"
	"github.com/leoferamos/aroma-sense/internal/storage"
	"gorm.io/gorm"
)

// AppHandlers contains all initialized handlers
type AppHandlers struct {
	UserHandler              *userhandler.UserHandler
	AdminUserHandler         *admin.AdminUserHandler
	ProductHandler           *product.ProductHandler
	CartHandler              *carthandler.CartHandler
	OrderHandler             *orderhandler.OrderHandler
	PasswordResetHandler     *auth.PasswordResetHandler
	ShippingHandler          *shipping.ShippingHandler
	ReviewHandler            *reviewhandler.ReviewHandler
	AIHandler                *aihandler.AIHandler
	ChatHandler              *chathandler.ChatHandler
	AuditLogHandler          *loghandler.AuditLogHandler
	AdminContestationHandler *admin.AdminContestationHandler
	AdminReviewReportHandler *admin.AdminReviewReportHandler
	PaymentHandler           *paymenthandler.PaymentHandler
}

// AppServices contains service instances needed for jobs
type AppServices struct {
	AdminUserService serviceadmin.AdminUserService
	AuditLogService  servicelog.AuditLogService
	LgpdService      servicelgpd.LgpdService
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
