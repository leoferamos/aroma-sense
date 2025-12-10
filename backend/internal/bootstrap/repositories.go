package bootstrap

import (
	"github.com/leoferamos/aroma-sense/internal/repository"
	"gorm.io/gorm"
)

// repositories holds all repository instances
type repositories struct {
	user             repository.UserRepository
	product          repository.ProductRepository
	cart             repository.CartRepository
	order            repository.OrderRepository
	payment          repository.PaymentRepository
	resetToken       repository.ResetTokenRepository
	review           repository.ReviewRepository
	reviewReport     repository.ReviewReportRepository
	auditLog         repository.AuditLogRepository
	userContestation repository.UserContestationRepository
}

// initializeRepositories creates all repository instances
func initializeRepositories(db *gorm.DB) *repositories {
	return &repositories{
		user:             repository.NewUserRepository(db),
		product:          repository.NewProductRepository(db),
		cart:             repository.NewCartRepository(db),
		order:            repository.NewOrderRepository(db),
		payment:          repository.NewPaymentRepository(db),
		resetToken:       repository.NewResetTokenRepository(db),
		review:           repository.NewReviewRepository(db),
		reviewReport:     repository.NewReviewReportRepository(db),
		auditLog:         repository.NewAuditLogRepository(db),
		userContestation: repository.NewUserContestationRepository(db),
	}
}
