package service

import (
	"strings"

	"github.com/leoferamos/aroma-sense/internal/apperror"
	"github.com/leoferamos/aroma-sense/internal/model"
	"github.com/leoferamos/aroma-sense/internal/repository"
	"github.com/leoferamos/aroma-sense/internal/validation"
	"golang.org/x/crypto/bcrypt"
)

// UserProfileService defines the interface for user profile-related business logic
type UserProfileService interface {
	GetByPublicID(publicID string) (*model.User, error)
	UpdateDisplayName(publicID string, displayName string) (*model.User, error)
	SetPasswordHash(publicID string, hashedPassword string) error
	ChangePassword(publicID string, currentPassword string, newPassword string) error
}

type userProfileService struct {
	repo            repository.UserRepository
	auditLogService AuditLogService
}

func NewUserProfileService(repo repository.UserRepository, auditLogService AuditLogService) UserProfileService {
	return &userProfileService{repo: repo, auditLogService: auditLogService}
}

// GetByPublicID returns the user by public id
func (s *userProfileService) GetByPublicID(publicID string) (*model.User, error) {
	if publicID == "" {
		return nil, apperror.NewCodeMessage("unauthenticated", "unauthenticated")
	}
	return s.repo.FindByPublicID(publicID)
}

// UpdateDisplayName updates the user's display name with validation
func (s *userProfileService) UpdateDisplayName(publicID string, displayName string) (*model.User, error) {
	if publicID == "" {
		return nil, apperror.NewCodeMessage("unauthenticated", "unauthenticated")
	}
	trimmed := strings.TrimSpace(displayName)
	if len(trimmed) < 2 {
		return nil, apperror.NewCodeMessage("display_name_too_short", "display_name too short")
	}
	if len(trimmed) > 50 {
		return nil, apperror.NewCodeMessage("display_name_too_long", "display_name too long")
	}
	user, err := s.repo.FindByPublicID(publicID)
	if err != nil {
		return nil, err
	}

	// Store old values for audit log
	oldUser := *user

	dn := strings.TrimSpace(displayName)
	user.DisplayName = &dn
	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	// Log display name update
	if s.auditLogService != nil {
		s.auditLogService.LogUserUpdate(user.ID, user.ID, &oldUser, user)
	}

	return user, nil
}

// SetPasswordHash updates a user's password hash (low-level method)
func (s *userProfileService) SetPasswordHash(publicID string, hashedPassword string) error {
	user, err := s.repo.FindByPublicID(publicID)
	if err != nil {
		return err
	}

	// Store old values for audit log
	oldUser := *user

	user.PasswordHash = hashedPassword
	if err := s.repo.Update(user); err != nil {
		return err
	}

	// Log password change
	if s.auditLogService != nil {
		s.auditLogService.LogUserUpdate(user.ID, user.ID, &oldUser, user)
	}

	return nil
}

// ChangePassword changes the user's password after verifying the current password
func (s *userProfileService) ChangePassword(publicID string, currentPassword string, newPassword string) error {
	user, err := s.repo.FindByPublicID(publicID)
	if err != nil {
		return err
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword)); err != nil {
		return apperror.NewCodeMessage("current_password_incorrect", "current password is incorrect")
	}

	// Validate new password
	if err := validation.ValidatePassword(newPassword, user.Email); err != nil {
		return err
	}

	// Check if new password is the same as current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(newPassword)); err == nil {
		return apperror.NewCodeMessage("new_password_same", "new password must be different from current password")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return apperror.NewCodeMessage("password_process_failed", "failed to process password")
	}

	// Update password
	return s.SetPasswordHash(publicID, string(hashedPassword))
}
