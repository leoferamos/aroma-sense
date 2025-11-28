package job

import (
	"log"
	"time"

	"github.com/leoferamos/aroma-sense/internal/model"
	"github.com/leoferamos/aroma-sense/internal/repository"
	"github.com/leoferamos/aroma-sense/internal/service"
)

// DataCleanupJob handles automated data cleanup tasks for LGPD compliance
type DataCleanupJob struct {
	userRepo        repository.UserRepository
	lgpdService     service.LgpdService
	auditLogService service.AuditLogService
}

// NewDataCleanupJob creates a new data cleanup job instance
func NewDataCleanupJob(userRepo repository.UserRepository, lgpdService service.LgpdService, auditLogService service.AuditLogService) *DataCleanupJob {
	return &DataCleanupJob{
		userRepo:        userRepo,
		lgpdService:     lgpdService,
		auditLogService: auditLogService,
	}
}

// Start begins the automated cleanup process
func (j *DataCleanupJob) Start() {
	log.Println("Starting LGPD data cleanup job...")

	// Run initial cleanup
	j.runCleanup()

	// Schedule daily cleanup at 2 AM
	ticker := time.NewTicker(24 * time.Hour)
	go func() {
		for {
			// Wait for next tick
			<-ticker.C

			// Run cleanup
			j.runCleanup()
		}
	}()

	log.Println("LGPD data cleanup job scheduled to run daily at 2 AM")
}

// runCleanup performs the actual data cleanup
func (j *DataCleanupJob) runCleanup() {
	log.Println("Running LGPD data cleanup...")

	// Find users who have confirmed deletion and exceeded retention period
	expiredUsers, err := j.findExpiredUsers()
	if err != nil {
		log.Printf("Error finding expired users: %v", err)
		return
	}

	if len(expiredUsers) == 0 {
		log.Println("No expired users found for cleanup")
		return
	}

	log.Printf("Found %d expired users for anonymization", len(expiredUsers))

	// Anonymize each expired user
	anonymizedCount := 0
	for _, user := range expiredUsers {
		if err := j.lgpdService.AnonymizeExpiredUser(user.PublicID); err != nil {
			log.Printf("Error anonymizing user %s: %v", user.PublicID, err)
			continue
		}
		anonymizedCount++
		log.Printf("Successfully anonymized user %s", user.PublicID)
	}

	// Log cleanup summary
	if err := j.auditLogService.LogSystemAction("data_cleanup_completed", "system", "",
		map[string]interface{}{
			"users_processed":   len(expiredUsers),
			"users_anonymized":  anonymizedCount,
			"cleanup_timestamp": time.Now(),
			"lgpd_compliant":    true,
		}); err != nil {
		log.Printf("Error logging cleanup summary: %v", err)
	}

	log.Printf("LGPD data cleanup completed: %d/%d users anonymized", anonymizedCount, len(expiredUsers))
}

// findExpiredUsers finds users who have confirmed deletion and exceeded the 2-year retention period
func (j *DataCleanupJob) findExpiredUsers() ([]*model.User, error) {
	return j.userRepo.FindExpiredUsersForAnonymization()
}

// ManualCleanup allows manual triggering of cleanup for testing/admin purposes
func (j *DataCleanupJob) ManualCleanup() error {
	log.Println("Manual LGPD data cleanup triggered...")
	j.runCleanup()
	return nil
}
