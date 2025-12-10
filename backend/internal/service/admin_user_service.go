package service

import (
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/leoferamos/aroma-sense/internal/apperror"
	"github.com/leoferamos/aroma-sense/internal/model"
	"github.com/leoferamos/aroma-sense/internal/notification"
	"github.com/leoferamos/aroma-sense/internal/repository"
	"github.com/leoferamos/aroma-sense/internal/utils"
	"github.com/leoferamos/aroma-sense/internal/validation"
)

// AdminUserService defines the interface for admin user management business logic
type AdminUserService interface {
	ListUsers(limit int, offset int, filters map[string]interface{}) ([]*model.User, int64, error)
	GetUserByID(id uint) (*model.User, error)
	UpdateUserRole(userID uint, newRole string, adminPublicID string) error
	CreateAdminUser(email string, password string, displayName string, superAdminPublicID string) (*model.User, error)
	DeactivateUser(userID uint, adminPublicID string, reason string, notes string, suspensionUntil *time.Time) error
	AdminReactivateUser(userID uint, adminPublicID string, reason string) error
}

type adminUserService struct {
	repo            repository.UserRepository
	auditLogService AuditLogService
	notifier        notification.NotificationService
}

func NewAdminUserService(repo repository.UserRepository, auditLogService AuditLogService, notifier notification.NotificationService) AdminUserService {
	return &adminUserService{repo: repo, auditLogService: auditLogService, notifier: notifier}
}

// CreateAdminUser allows super admin to create a new admin user
func (s *adminUserService) CreateAdminUser(email string, password string, displayName string, superAdminPublicID string) (*model.User, error) {
	if err := validation.ValidatePassword(password, email); err != nil {
		return nil, err
	}
	if existing, err := s.repo.FindByEmail(email); err == nil && existing != nil {
		return nil, apperror.NewCodeMessage("email_already_registered", "email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperror.NewCodeMessage("password_hash_failed", "failed to hash password")
	}

	user := model.User{
		Email:        strings.ToLower(email),
		PasswordHash: string(hashedPassword),
		Role:         "admin",
	}
	if displayName != "" {
		user.DisplayName = &displayName
	}

	if err := s.repo.Create(&user); err != nil {
		return nil, err
	}

	if s.auditLogService != nil {
		admin, err := s.repo.FindByPublicID(superAdminPublicID)
		if err == nil {
			s.auditLogService.LogAdminAction(admin.ID, user.ID, model.AuditActionRoleChanged, map[string]interface{}{
				"action":     "create_admin",
				"new_role":   "admin",
				"email":      utils.MaskEmail(user.Email),
				"created_by": superAdminPublicID,
			})
		}
	}

	return &user, nil
}

// ListUsers returns paginated list of users for admin (LGPD compliance)
func (s *adminUserService) ListUsers(limit int, offset int, filters map[string]interface{}) ([]*model.User, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	return s.repo.ListUsers(limit, offset, filters)
}

// GetUserByID returns user by database ID for admin
func (s *adminUserService) GetUserByID(id uint) (*model.User, error) {
	return s.repo.FindByID(id)
}

// UpdateUserRole updates user role with admin tracking
func (s *adminUserService) UpdateUserRole(userID uint, newRole string, adminPublicID string) error {
	// Validate role
	if newRole != "admin" && newRole != "client" {
		return apperror.NewCodeMessage("invalid_role", "invalid role")
	}

	// Prevent admin from changing their own role
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return err
	}
	if user.PublicID == adminPublicID {
		return apperror.NewCodeMessage("cannot_change_own_role", "cannot change your own role")
	}

	// Get admin user for audit log
	admin, err := s.repo.FindByPublicID(adminPublicID)
	if err != nil {
		return err
	}

	oldRole := user.Role
	if err := s.repo.UpdateRole(userID, newRole); err != nil {
		return err
	}

	// Log role change
	if s.auditLogService != nil {
		s.auditLogService.LogAdminAction(admin.ID, userID, model.AuditActionRoleChanged,
			map[string]interface{}{
				"old_role": oldRole,
				"new_role": newRole,
			})
	}

	return nil
}

// DeactivateUser soft deletes a user account with enhanced LGPD compliance
func (s *adminUserService) DeactivateUser(userID uint, adminPublicID string, reason string, notes string, suspensionUntil *time.Time) error {
	// Get admin user for audit log
	admin, err := s.repo.FindByPublicID(adminPublicID)
	if err != nil {
		return err
	}

	now := time.Now()
	if suspensionUntil != nil && suspensionUntil.Before(now) {
		return apperror.NewCodeMessage("suspension_until_past", "suspensionUntil cannot be in the past")
	}
	if err := s.repo.DeactivateUser(userID, adminPublicID, now, reason, notes, suspensionUntil); err != nil {
		return err
	}

	// Log user deactivation
	if s.auditLogService != nil {
		s.auditLogService.LogAdminAction(admin.ID, userID, model.AuditActionUserDeactivated,
			map[string]interface{}{
				"reason":           reason,
				"notes":            notes,
				"suspension_until": suspensionUntil,
				"deactivated_at":   now,
			})
	}

	// Send deactivation notification email
	user, err := s.repo.FindByID(userID)
	if err == nil && s.notifier != nil {
		deadlineStr := "7 dias a partir da desativação"
		if user.ContestationDeadline != nil {
			deadlineStr = user.ContestationDeadline.Format("02/01/2006 15:04")
		}
		s.notifier.SendAccountDeactivated(user.Email, reason, deadlineStr)
	}

	// Invalidate refresh token
	if user != nil {
		_ = s.repo.UpdateRefreshToken(user.ID, nil, nil)
	}

	return nil
}

// AdminReactivateUser allows admin to reactivate a user account after contestation review
func (s *adminUserService) AdminReactivateUser(userID uint, adminPublicID string, reason string) error {
	// Get admin user for audit log
	admin, err := s.repo.FindByPublicID(adminPublicID)
	if err != nil {
		return err
	}

	user, err := s.repo.FindByID(userID)
	if err != nil {
		return err
	}

	// Check if user is deactivated
	if user.DeactivatedAt == nil {
		return apperror.NewCodeMessage("user_not_deactivated", "user is not deactivated")
	}

	// Reactivate user
	user.DeactivatedAt = nil
	user.DeactivatedBy = nil
	user.DeactivationReason = nil
	user.DeactivationNotes = nil
	user.SuspensionUntil = nil
	user.ReactivationRequested = false
	user.ContestationDeadline = nil

	if err := s.repo.Update(user); err != nil {
		return err
	}

	// Log reactivation
	if s.auditLogService != nil {
		s.auditLogService.LogAdminAction(admin.ID, userID, model.AuditActionUserReactivated,
			map[string]interface{}{
				"reason": reason,
			})
	}

	// Send reactivation result email
	if s.notifier != nil {
		s.notifier.SendContestationResult(user.Email, true, reason)
	}

	return nil
}
