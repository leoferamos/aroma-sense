package job

import (
	"log"
	"time"

	"github.com/leoferamos/aroma-sense/internal/model"
	"github.com/leoferamos/aroma-sense/internal/repository"
	"github.com/leoferamos/aroma-sense/internal/service"
)

// AutoConfirmJob handles automatic confirmation of account deletion after the cooling-off period
type AutoConfirmJob struct {
	userRepo        repository.UserRepository
	lgpdService     service.LgpdService
	auditLogService service.AuditLogService
}

// NewAutoConfirmJob creates a new auto-confirm job instance
func NewAutoConfirmJob(userRepo repository.UserRepository, lgpdService service.LgpdService, auditLogService service.AuditLogService) *AutoConfirmJob {
	return &AutoConfirmJob{
		userRepo:        userRepo,
		lgpdService:     lgpdService,
		auditLogService: auditLogService,
	}
}

// Start schedules daily auto-confirm runs
func (j *AutoConfirmJob) Start() {
	log.Println("Starting account auto-confirm job...")

	// Run initial pass
	j.runAutoConfirm()

	// Schedule daily run
	ticker := time.NewTicker(24 * time.Hour)
	go func() {
		for {
			<-ticker.C
			j.runAutoConfirm()
		}
	}()

	log.Println("Account auto-confirm job scheduled to run daily")
}

// runAutoConfirm performs the actual auto-confirm work
func (j *AutoConfirmJob) runAutoConfirm() {
	log.Println("Running account auto-confirm job...")

	// cutoff is 7 days ago
	cutoff := time.Now().Add(-7 * 24 * time.Hour)

	candidates, err := j.userRepo.FindUsersPendingAutoConfirm(cutoff)
	if err != nil {
		log.Printf("Error querying pending auto-confirm users: %v", err)
		return
	}

	if len(candidates) == 0 {
		log.Println("No users pending auto-confirm")
		return
	}

	processed := 0
	skipped := 0

	for _, u := range candidates {
		// Re-check active dependencies before confirming
		hasDeps, err := j.userRepo.HasActiveDependencies(u.PublicID)
		if err != nil {
			log.Printf("Error checking dependencies for user %s: %v", u.PublicID, err)
			skipped++
			continue
		}
		if hasDeps {
			// Log and skip - team may need to review
			if j.auditLogService != nil {
				j.auditLogService.LogSystemAction("auto_confirm_skipped_active_dependencies", "system", u.PublicID,
					map[string]interface{}{"reason": "active_dependencies"})
			}
			log.Printf("Skipping auto-confirm for user %s due to active dependencies", u.PublicID)
			skipped++
			continue
		}

		// Attempt to confirm deletion via LGPD service (this will validate cooling-off has passed)
		if err := j.lgpdService.ConfirmAccountDeletion(u.PublicID); err != nil {
			log.Printf("Error confirming deletion for user %s: %v", u.PublicID, err)
			skipped++
			continue
		}

		// Log system audit
		if j.auditLogService != nil {
			j.auditLogService.LogSystemAction(model.AuditActionDeletionConfirmed, "system", u.PublicID,
				map[string]interface{}{"auto_confirm": true, "confirmed_at": time.Now()})
		}

		processed++
		log.Printf("Auto-confirmed deletion for user %s", u.PublicID)
	}

	log.Printf("Account auto-confirm completed: %d processed, %d skipped", processed, skipped)
}

// ManualRun allows manual triggering of the auto-confirm job (for testing/admin purposes)
func (j *AutoConfirmJob) ManualRun() error {
	log.Println("Manual auto-confirm run triggered...")
	j.runAutoConfirm()
	return nil
}
