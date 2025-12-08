package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/leoferamos/aroma-sense/internal/dto"
	"github.com/leoferamos/aroma-sense/internal/model"
	"github.com/leoferamos/aroma-sense/internal/notification"
	"github.com/leoferamos/aroma-sense/internal/repository"
)

// LgpdService defines the interface for LGPD/GDPR compliance business logic
type LgpdService interface {
	ExportUserData(publicID string) (*dto.UserExportResponse, error)
	RequestAccountDeletion(publicID string) error
	ConfirmAccountDeletion(publicID string) error
	CancelAccountDeletion(publicID string) error
	AnonymizeExpiredUser(publicID string) error
	RequestContestation(publicID string, reason string) error
	ProcessPendingDeletions() error
	ProcessExpiredAnonymizations() error
}

type lgpdService struct {
	repo             repository.UserRepository
	userContestation repository.UserContestationRepository
	auditLogService  AuditLogService
	notifier         notification.NotificationService
}

func NewLgpdService(repo repository.UserRepository, userContestationRepo repository.UserContestationRepository, auditLogService AuditLogService, notifier notification.NotificationService) LgpdService {
	return &lgpdService{repo: repo, userContestation: userContestationRepo, auditLogService: auditLogService, notifier: notifier}
}

// ExportUserData exports all user data for GDPR compliance
func (s *lgpdService) ExportUserData(publicID string) (*dto.UserExportResponse, error) {
	user, err := s.repo.FindByPublicID(publicID)
	if err != nil {
		return nil, err
	}

	return &dto.UserExportResponse{
		PublicID:            user.PublicID,
		Email:               user.Email,
		Role:                user.Role,
		DisplayName:         user.DisplayName,
		CreatedAt:           user.CreatedAt,
		LastLoginAt:         user.LastLoginAt,
		DeactivatedAt:       user.DeactivatedAt,
		DeletionRequestedAt: user.DeletionRequestedAt,
		DeletionConfirmedAt: user.DeletionConfirmedAt,
	}, nil
}

// RequestAccountDeletion initiates account deletion process with 7-day cooling off period (LGPD compliance)
func (s *lgpdService) RequestAccountDeletion(publicID string) error {
	if publicID == "" {
		return errors.New("unauthenticated")
	}

	// Check if user has active dependencies
	hasDependencies, err := s.repo.HasActiveDependencies(publicID)
	if err != nil {
		return errors.New("failed to check account dependencies")
	}
	if hasDependencies {
		return errors.New("cannot delete account with active orders - please cancel or complete all orders first")
	}

	// Check if already requested deletion
	user, err := s.repo.FindByPublicID(publicID)
	if err != nil {
		return err
	}
	if user.DeletionRequestedAt != nil {
		return errors.New("account deletion already requested")
	}

	now := time.Now()
	if err := s.repo.RequestAccountDeletion(publicID, now); err != nil {
		return err
	}

	// Send deletion requested email
	if s.notifier != nil {
		_ = s.notifier.SendDeletionRequested(user.Email, "")
	}

	// Log deletion request
	if s.auditLogService != nil {
		s.auditLogService.LogDeletionAction(nil, user.ID, model.AuditActionDeletionRequested,
			map[string]interface{}{
				"cooling_off_period_days": 7,
				"retention_period_years":  2,
			})
	}

	return nil
}

// ConfirmAccountDeletion confirms account deletion after cooling off period (LGPD compliance)
func (s *lgpdService) ConfirmAccountDeletion(publicID string) error {
	if publicID == "" {
		return errors.New("unauthenticated")
	}

	user, err := s.repo.FindByPublicID(publicID)
	if err != nil {
		return err
	}

	// Check if deletion was requested
	if user.DeletionRequestedAt == nil {
		return errors.New("no deletion request found")
	}

	// Check if cooling off period (7 days) has passed
	coolingOffPeriod := user.DeletionRequestedAt.Add(7 * 24 * time.Hour)
	if time.Now().Before(coolingOffPeriod) {
		return errors.New("cooling off period not yet expired - please wait 7 days from request")
	}

	now := time.Now()
	if err := s.repo.ConfirmAccountDeletion(publicID, now); err != nil {
		return err
	}

	// Log deletion confirmation
	if s.auditLogService != nil {
		s.auditLogService.LogDeletionAction(nil, user.ID, model.AuditActionUserDeleted,
			map[string]interface{}{
				"cooling_off_expired":    true,
				"retention_period_years": 2,
				"deletion_requested_at":  user.DeletionRequestedAt,
			})
	}

	// Send deletion confirmed email
	if s.notifier != nil {
		_ = s.notifier.SendDeletionAutoConfirmed(user.Email)
	}

	return nil
}

// CancelAccountDeletion cancels a pending account deletion request
func (s *lgpdService) CancelAccountDeletion(publicID string) error {
	if publicID == "" {
		return errors.New("unauthenticated")
	}

	user, err := s.repo.FindByPublicID(publicID)
	if err != nil {
		return err
	}

	if user.DeletionRequestedAt == nil {
		return errors.New("no deletion request to cancel")
	}

	// capture requested time before clearing
	requestedAt := user.DeletionRequestedAt

	// Clear deletion request
	user.DeletionRequestedAt = nil
	if err := s.repo.Update(user); err != nil {
		return err
	}

	// Log deletion cancellation
	if s.auditLogService != nil {
		s.auditLogService.LogDeletionAction(nil, user.ID, model.AuditActionDeletionCancelled,
			map[string]interface{}{
				"reason":                "user_cancelled",
				"deletion_requested_at": requestedAt,
			})
	}

	// Send cancellation email
	if s.notifier != nil {
		_ = s.notifier.SendDeletionCancelled(user.Email)
	}

	return nil
}

// AnonymizeExpiredUser anonymizes user data after retention period (LGPD compliance)
func (s *lgpdService) AnonymizeExpiredUser(publicID string) error {
	user, err := s.repo.FindByPublicID(publicID)
	if err != nil {
		return err
	}

	// Verify deletion was confirmed and retention period has passed (2 years)
	if user.DeletionConfirmedAt == nil {
		return errors.New("user has not confirmed account deletion")
	}

	retentionPeriod := user.DeletionConfirmedAt.Add(5 * 365 * 24 * time.Hour) // 5 years
	if time.Now().Before(retentionPeriod) {
		return errors.New("retention period not yet expired")
	}

	// Anonymize personal data while keeping necessary records for compliance
	anonymizedEmail := fmt.Sprintf("deleted-%s@anonymous.local", user.PublicID[:8])
	anonymizedDisplayName := "Usuário Excluído"

	if err := s.repo.AnonymizeUser(publicID, anonymizedEmail, anonymizedDisplayName); err != nil {
		return err
	}

	// Log data anonymization
	if s.auditLogService != nil {
		s.auditLogService.LogSystemAction(model.AuditActionDataAnonymized, "user", publicID,
			map[string]interface{}{
				"retention_period_expired": true,
				"deletion_confirmed_at":    user.DeletionConfirmedAt,
				"anonymized_email":         anonymizedEmail,
				"anonymized_display_name":  anonymizedDisplayName,
				"lgpd_compliant":           true,
			})
	}

	// Send data anonymized email
	if s.notifier != nil {
		// send to previous email address
		_ = s.notifier.SendDataAnonymized(user.Email)
	}

	return nil
}

// RequestContestation allows user to contest account deactivation (LGPD compliance)
func (s *lgpdService) RequestContestation(publicID string, reason string) error {
	if publicID == "" {
		return errors.New("unauthenticated")
	}

	user, err := s.repo.FindByPublicID(publicID)
	if err != nil {
		return err
	}

	// Check if user is deactivated
	if user.DeactivatedAt == nil {
		return errors.New("account is not deactivated")
	}

	// Check if contestation deadline has passed
	if user.ContestationDeadline != nil && user.ContestationDeadline.Before(time.Now()) {
		return errors.New("contestation deadline has expired")
	}

	// Check if already requested reactivation
	if user.ReactivationRequested {
		return errors.New("reactivation already requested")
	}

	// Set contestation deadline if not set (7 days from deactivation)
	if user.ContestationDeadline == nil {
		deadline := user.DeactivatedAt.Add(7 * 24 * time.Hour)
		user.ContestationDeadline = &deadline
	}

	// Create contestation record
	contest := &model.UserContestation{
		UserID:      user.ID,
		Reason:      reason,
		Status:      "pending",
		RequestedAt: time.Now(),
	}
	if err := s.userContestation.Create(contest); err != nil {
		return err
	}

	// Log contestation request
	if s.auditLogService != nil {
		s.auditLogService.LogUserAction(nil, &user.ID, model.AuditActionUserReactivated,
			map[string]interface{}{
				"action_type": "contestation_requested",
				"reason":      reason,
				"deadline":    user.ContestationDeadline,
			})
	}

	// Send contestation received confirmation email
	if s.notifier != nil {
		_ = s.notifier.SendContestationReceived(user.Email)
	}

	return nil
}

// ProcessPendingDeletions automatically confirms deletions after 7 days (daily job)
func (s *lgpdService) ProcessPendingDeletions() error {
	cutoff := time.Now().Add(-7 * 24 * time.Hour) // 7 days ago
	// Find users with deletion_requested_at > 7 days ago and not yet confirmed
	users, err := s.repo.FindUsersPendingAutoConfirm(cutoff)
	if err != nil {
		return fmt.Errorf("failed to find pending deletions: %w", err)
	}

	for _, user := range users {
		if err := s.ConfirmAccountDeletion(user.PublicID); err != nil {
			// Log error but continue processing others
			fmt.Printf("Failed to confirm deletion for user %s: %v\n", user.PublicID, err)
		}
	}
	return nil
}

// ProcessExpiredAnonymizations anonymizes users after 5 years (daily job)
func (s *lgpdService) ProcessExpiredAnonymizations() error {
	// Find users with deletion_confirmed_at > 5 years ago and not yet anonymized
	users, err := s.repo.FindExpiredUsersForAnonymization()
	if err != nil {
		return fmt.Errorf("failed to find users for anonymization: %w", err)
	}

	for _, user := range users {
		if err := s.AnonymizeExpiredUser(user.PublicID); err != nil {
			// Log error but continue
			fmt.Printf("Failed to anonymize user %s: %v\n", user.PublicID, err)
		}
	}
	return nil
}
