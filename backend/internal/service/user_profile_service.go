package service

import (
	"errors"
	"strings"

	"github.com/leoferamos/aroma-sense/internal/model"
	"github.com/leoferamos/aroma-sense/internal/repository"
)

// UserProfileService defines the interface for user profile-related business logic
type UserProfileService interface {
	GetByPublicID(publicID string) (*model.User, error)
	UpdateDisplayName(publicID string, displayName string) (*model.User, error)
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
		return nil, errors.New("unauthenticated")
	}
	return s.repo.FindByPublicID(publicID)
}

// UpdateDisplayName updates the user's display name with validation
func (s *userProfileService) UpdateDisplayName(publicID string, displayName string) (*model.User, error) {
	if publicID == "" {
		return nil, errors.New("unauthenticated")
	}
	trimmed := strings.TrimSpace(displayName)
	if len(trimmed) < 2 {
		return nil, errors.New("display_name too short")
	}
	if len(trimmed) > 50 {
		return nil, errors.New("display_name too long")
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
