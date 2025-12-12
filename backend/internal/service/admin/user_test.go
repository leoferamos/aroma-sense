package service

import (
	"errors"
	"testing"
	"time"

	"github.com/leoferamos/aroma-sense/internal/dto"
	"github.com/leoferamos/aroma-sense/internal/model"
	"github.com/leoferamos/aroma-sense/internal/notification"
	"github.com/leoferamos/aroma-sense/internal/repository"
	logservice "github.com/leoferamos/aroma-sense/internal/service/log"
	"github.com/stretchr/testify/assert"
)

// Mock implementations
type mockUserRepo struct {
	createErr                  error
	findByEmailUser            *model.User
	findByEmailErr             error
	findByRefreshTokenHashUser *model.User
	findByRefreshTokenHashErr  error
	findByPublicIDUser         *model.User
	findByPublicIDErr          error
	updateErr                  error
	updateRoleErr              error
	listUsers                  []*model.User
	listCount                  int64
	listErr                    error
	findByIDUser               *model.User
	findByIDErr                error
	deactivateUserErr          error
	updateRefreshTokenErr      error
}

func (m *mockUserRepo) Create(user *model.User) error {
	return m.createErr
}

func (m *mockUserRepo) FindByEmail(email string) (*model.User, error) {
	return m.findByEmailUser, m.findByEmailErr
}

func (m *mockUserRepo) FindByRefreshTokenHash(hash string) (*model.User, error) {
	return nil, nil
}

func (m *mockUserRepo) FindByPublicID(publicID string) (*model.User, error) {
	return m.findByPublicIDUser, m.findByPublicIDErr
}

func (m *mockUserRepo) Update(user *model.User) error {
	return m.updateErr
}

func (m *mockUserRepo) UpdateRefreshToken(userID uint, hash *string, expiresAt *time.Time) error {
	return m.updateRefreshTokenErr
}

func (m *mockUserRepo) ListUsers(limit int, offset int, filters map[string]interface{}) ([]*model.User, int64, error) {
	return m.listUsers, m.listCount, m.listErr
}

func (m *mockUserRepo) FindByID(id uint) (*model.User, error) {
	return m.findByIDUser, m.findByIDErr
}

func (m *mockUserRepo) UpdateRole(userID uint, newRole string) error {
	return m.updateRoleErr
}

func (m *mockUserRepo) DeactivateUser(userID uint, adminPublicID string, deactivatedAt time.Time, reason string, notes string, suspensionUntil *time.Time) error {
	return m.deactivateUserErr
}

func (m *mockUserRepo) RequestAccountDeletion(publicID string, requestedAt time.Time) error {
	return nil
}

func (m *mockUserRepo) ConfirmAccountDeletion(publicID string, confirmedAt time.Time) error {
	return nil
}

func (m *mockUserRepo) HasActiveDependencies(publicID string) (bool, error) {
	return false, nil
}

func (m *mockUserRepo) AnonymizeUser(publicID string, anonymizedEmail string, anonymizedDisplayName string) error {
	return nil
}

func (m *mockUserRepo) FindExpiredUsersForAnonymization() ([]*model.User, error) {
	return nil, nil
}

func (m *mockUserRepo) FindUsersPendingAutoConfirm(cutoff time.Time) ([]*model.User, error) {
	return nil, nil
}

func (m *mockUserRepo) DeleteByPublicID(publicID string) error {
	return nil
}

type mockAuditLogService struct {
	logAdminActionErr error
}

func (m *mockAuditLogService) LogUserAction(actorID *uint, userID *uint, action model.AuditAction, details map[string]interface{}) error {
	return nil
}

func (m *mockAuditLogService) LogUserUpdate(actorID uint, userID uint, oldUser, newUser *model.User) error {
	return nil
}

func (m *mockAuditLogService) LogAdminAction(adminID uint, userID uint, action model.AuditAction, details map[string]interface{}) error {
	return m.logAdminActionErr
}

func (m *mockAuditLogService) LogSystemAction(action model.AuditAction, resource, resourceID string, details map[string]interface{}) error {
	return nil
}

func (m *mockAuditLogService) LogDataAccess(userID uint, resource string, resourceID string) error {
	return nil
}

func (m *mockAuditLogService) LogDeletionAction(actorID *uint, userID uint, action model.AuditAction, details map[string]interface{}) error {
	return nil
}

func (m *mockAuditLogService) ListAuditLogs(filter *model.AuditLogFilter) ([]*model.AuditLog, int64, error) {
	return nil, 0, nil
}

func (m *mockAuditLogService) GetAuditLogByID(id uint) (*model.AuditLog, error) {
	return nil, nil
}

func (m *mockAuditLogService) GetUserAuditLogs(userID uint, limit, offset int) ([]*model.AuditLog, int64, error) {
	return nil, 0, nil
}

func (m *mockAuditLogService) GetResourceAuditLogs(resource, resourceID string) ([]*model.AuditLog, error) {
	return nil, nil
}

func (m *mockAuditLogService) GetAuditSummary(startDate, endDate *time.Time) (*model.AuditLogSummary, error) {
	return nil, nil
}

func (m *mockAuditLogService) CleanupOldLogs(retentionDays int) error {
	return nil
}

func (m *mockAuditLogService) ConvertAuditLogToResponse(auditLog *model.AuditLog) dto.AuditLogResponse {
	return dto.AuditLogResponse{}
}

func (m *mockAuditLogService) ConvertAuditLogToResponseDetailed(auditLog *model.AuditLog) dto.AuditLogResponse {
	return dto.AuditLogResponse{}
}

func (m *mockAuditLogService) ConvertAuditLogsToResponse(auditLogs []*model.AuditLog) []dto.AuditLogResponse {
	return []dto.AuditLogResponse{}
}

func (m *mockAuditLogService) ConvertAuditLogSummaryToResponse(summary *model.AuditLogSummary) dto.AuditLogSummaryResponse {
	return dto.AuditLogSummaryResponse{}
}

type mockNotificationService struct {
	sendAccountDeactivatedErr error
	sendContestationResultErr error
}

func (m *mockNotificationService) SendPasswordResetCode(to, code string) error {
	return nil
}

func (m *mockNotificationService) SendWelcomeEmail(to, name string) error {
	return nil
}

func (m *mockNotificationService) SendOrderConfirmation(to string, order *model.Order) error {
	return nil
}

func (m *mockNotificationService) SendAccountDeactivated(to, reason string, contestationDeadline string) error {
	return m.sendAccountDeactivatedErr
}

func (m *mockNotificationService) SendContestationReceived(to string) error {
	return nil
}

func (m *mockNotificationService) SendContestationResult(to string, approved bool, reason string) error {
	return m.sendContestationResultErr
}

func (m *mockNotificationService) SendDeletionRequested(to string, cancelLink string) error {
	return nil
}

func (m *mockNotificationService) SendDeletionAutoConfirmed(to string) error {
	return nil
}

func (m *mockNotificationService) SendDeletionCancelled(to string) error {
	return nil
}

func (m *mockNotificationService) SendDataAnonymized(to string) error {
	return nil
}

func (m *mockNotificationService) SendPromotional(to, subject, htmlBody string) error {
	return nil
}

// Test helpers
func createTestUser() *model.User {
	return &model.User{
		ID:       1,
		PublicID: "user123",
		Email:    "test@example.com",
		Role:     "client",
	}
}

func createTestAdmin() *model.User {
	return &model.User{
		ID:       2,
		PublicID: "admin123",
		Email:    "admin@example.com",
		Role:     "admin",
	}
}

// Tests
func TestAdminUserService_CreateAdminUser(t *testing.T) {
	tests := []struct {
		name               string
		email              string
		password           string
		displayName        string
		superAdminPublicID string
		mockSetup          func() (repository.UserRepository, logservice.AuditLogService, notification.NotificationService)
		expectedErr        bool
		expectedErrCode    string
	}{
		{
			name:               "success",
			email:              "newadmin@example.com",
			password:           "ValidPass123!",
			displayName:        "New Admin",
			superAdminPublicID: "admin123",
			mockSetup: func() (repository.UserRepository, logservice.AuditLogService, notification.NotificationService) {
				userRepo := &mockUserRepo{
					findByEmailErr: errors.New("not found"),
					findByPublicIDUser: &model.User{
						ID:       1,
						PublicID: "admin123",
						Email:    "admin@example.com",
					},
				}
				auditSvc := &mockAuditLogService{}
				notifier := &mockNotificationService{}
				return userRepo, auditSvc, notifier
			},
			expectedErr: false,
		},
		{
			name:               "email already registered",
			email:              "existing@example.com",
			password:           "ValidPass123!",
			displayName:        "",
			superAdminPublicID: "admin123",
			mockSetup: func() (repository.UserRepository, logservice.AuditLogService, notification.NotificationService) {
				userRepo := &mockUserRepo{
					findByEmailUser: &model.User{Email: "existing@example.com"},
				}
				auditSvc := &mockAuditLogService{}
				notifier := &mockNotificationService{}
				return userRepo, auditSvc, notifier
			},
			expectedErr:     true,
			expectedErrCode: "email already registered",
		},
		{
			name:               "invalid password",
			email:              "newadmin@example.com",
			password:           "weak",
			displayName:        "",
			superAdminPublicID: "admin123",
			mockSetup: func() (repository.UserRepository, logservice.AuditLogService, notification.NotificationService) {
				userRepo := &mockUserRepo{
					findByEmailErr: errors.New("not found"),
				}
				auditSvc := &mockAuditLogService{}
				notifier := &mockNotificationService{}
				return userRepo, auditSvc, notifier
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo, auditSvc, notifier := tt.mockSetup()
			svc := NewAdminUserService(userRepo, auditSvc, notifier)

			user, err := svc.CreateAdminUser(tt.email, tt.password, tt.displayName, tt.superAdminPublicID)

			if tt.expectedErr {
				assert.Error(t, err)
				if tt.expectedErrCode != "" {
					assert.Contains(t, err.Error(), tt.expectedErrCode)
				}
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.email, user.Email)
				assert.Equal(t, "admin", user.Role)
			}
		})
	}
}

func TestAdminUserService_UpdateUserRole(t *testing.T) {
	tests := []struct {
		name            string
		userID          uint
		newRole         string
		adminPublicID   string
		mockSetup       func() (repository.UserRepository, logservice.AuditLogService, notification.NotificationService)
		expectedErr     bool
		expectedErrCode string
	}{
		{
			name:          "success update to admin",
			userID:        1,
			newRole:       "admin",
			adminPublicID: "admin123",
			mockSetup: func() (repository.UserRepository, logservice.AuditLogService, notification.NotificationService) {
				userRepo := &mockUserRepo{
					findByIDUser:       createTestUser(),
					findByPublicIDUser: createTestAdmin(),
				}
				auditSvc := &mockAuditLogService{}
				notifier := &mockNotificationService{}
				return userRepo, auditSvc, notifier
			},
			expectedErr: false,
		},
		{
			name:          "cannot change own role",
			userID:        2,
			newRole:       "client",
			adminPublicID: "admin123",
			mockSetup: func() (repository.UserRepository, logservice.AuditLogService, notification.NotificationService) {
				userRepo := &mockUserRepo{
					findByIDUser:       createTestAdmin(),
					findByPublicIDUser: createTestAdmin(),
				}
				auditSvc := &mockAuditLogService{}
				notifier := &mockNotificationService{}
				return userRepo, auditSvc, notifier
			},
			expectedErr:     true,
			expectedErrCode: "cannot change your own role",
		},
		{
			name:          "invalid role",
			userID:        1,
			newRole:       "invalid",
			adminPublicID: "admin123",
			mockSetup: func() (repository.UserRepository, logservice.AuditLogService, notification.NotificationService) {
				userRepo := &mockUserRepo{
					findByIDUser:       createTestUser(),
					findByPublicIDUser: createTestAdmin(),
				}
				auditSvc := &mockAuditLogService{}
				notifier := &mockNotificationService{}
				return userRepo, auditSvc, notifier
			},
			expectedErr:     true,
			expectedErrCode: "invalid role",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo, auditSvc, notifier := tt.mockSetup()
			svc := NewAdminUserService(userRepo, auditSvc, notifier)

			err := svc.UpdateUserRole(tt.userID, tt.newRole, tt.adminPublicID)

			if tt.expectedErr {
				assert.Error(t, err)
				if tt.expectedErrCode != "" {
					assert.Contains(t, err.Error(), tt.expectedErrCode)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAdminUserService_DeactivateUser(t *testing.T) {
	tests := []struct {
		name            string
		userID          uint
		adminPublicID   string
		reason          string
		notes           string
		suspensionUntil *time.Time
		mockSetup       func() (repository.UserRepository, logservice.AuditLogService, notification.NotificationService)
		expectedErr     bool
	}{
		{
			name:            "success",
			userID:          1,
			adminPublicID:   "admin123",
			reason:          "Violation of terms",
			notes:           "Repeated spam",
			suspensionUntil: nil,
			mockSetup: func() (repository.UserRepository, logservice.AuditLogService, notification.NotificationService) {
				userRepo := &mockUserRepo{
					findByPublicIDUser: createTestAdmin(),
					findByIDUser:       createTestUser(),
				}
				auditSvc := &mockAuditLogService{}
				notifier := &mockNotificationService{}
				return userRepo, auditSvc, notifier
			},
			expectedErr: false,
		},
		{
			name:            "suspension in past",
			userID:          1,
			adminPublicID:   "admin123",
			reason:          "Test",
			notes:           "",
			suspensionUntil: func() *time.Time { t := time.Now().Add(-time.Hour); return &t }(),
			mockSetup: func() (repository.UserRepository, logservice.AuditLogService, notification.NotificationService) {
				userRepo := &mockUserRepo{
					findByPublicIDUser: createTestAdmin(),
					findByIDUser:       createTestUser(),
				}
				auditSvc := &mockAuditLogService{}
				notifier := &mockNotificationService{}
				return userRepo, auditSvc, notifier
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo, auditSvc, notifier := tt.mockSetup()
			svc := NewAdminUserService(userRepo, auditSvc, notifier)

			err := svc.DeactivateUser(tt.userID, tt.adminPublicID, tt.reason, tt.notes, tt.suspensionUntil)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAdminUserService_AdminReactivateUser(t *testing.T) {
	tests := []struct {
		name          string
		userID        uint
		adminPublicID string
		reason        string
		mockSetup     func() (repository.UserRepository, logservice.AuditLogService, notification.NotificationService)
		expectedErr   bool
		expectedCode  string
	}{
		{
			name:          "success",
			userID:        1,
			adminPublicID: "admin123",
			reason:        "Appeal approved",
			mockSetup: func() (repository.UserRepository, logservice.AuditLogService, notification.NotificationService) {
				user := createTestUser()
				user.DeactivatedAt = &time.Time{}
				userRepo := &mockUserRepo{
					findByPublicIDUser: createTestAdmin(),
					findByIDUser:       user,
				}
				auditSvc := &mockAuditLogService{}
				notifier := &mockNotificationService{}
				return userRepo, auditSvc, notifier
			},
			expectedErr: false,
		},
		{
			name:          "user not deactivated",
			userID:        1,
			adminPublicID: "admin123",
			reason:        "Test",
			mockSetup: func() (repository.UserRepository, logservice.AuditLogService, notification.NotificationService) {
				userRepo := &mockUserRepo{
					findByPublicIDUser: createTestAdmin(),
					findByIDUser:       createTestUser(),
				}
				auditSvc := &mockAuditLogService{}
				notifier := &mockNotificationService{}
				return userRepo, auditSvc, notifier
			},
			expectedErr:  true,
			expectedCode: "user is not deactivated",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo, auditSvc, notifier := tt.mockSetup()
			svc := NewAdminUserService(userRepo, auditSvc, notifier)

			err := svc.AdminReactivateUser(tt.userID, tt.adminPublicID, tt.reason)

			if tt.expectedErr {
				assert.Error(t, err)
				if tt.expectedCode != "" {
					assert.Contains(t, err.Error(), tt.expectedCode)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
